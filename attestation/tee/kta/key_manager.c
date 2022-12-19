/*
kunpengsecl licensed under the Mulan PSL v2.
You can use this software according to the terms and conditions of
the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at:
    http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.

Author: leezhenxiang
Create: 2022-11-04
Description: key managing module in kta.
	1. 2022-11-04	leezhenxiang
		define the structures.
    2. 2022-12-01   waterh2o
        Implementation function sendrequest and getresponse
*/

#include <tee_defines.h>
#include <tee_object_api.h>
#include <tee_crypto_api.h>
#include <securec.h>
#include <kta_common.h>
#include <cJSON.h>

#define PARAM_COUNT 4
#define AES_KEY_SIZE 32
#define RSA_KEY_SIZE 256
#define NONCE_SIZE 12
#define TAG_SIZE 16

typedef struct ktaprivkey {
    uint8_t modulus[RSA_KEY_SIZE];
    uint8_t privateExponent[RSA_KEY_SIZE]
};

const TEE_UUID ktaUuid = {
    0x435dcafa, 0x0029, 0x4d53, { 0x97, 0xe8, 0xa7, 0xa1, 0x3a, 0x80, 0xc8, 0x2e }
};

const uint8_t exponent[3] = {0x01,0x00,0x01};

extern Cache cache;
extern CmdQueue cmdqueue;
extern ReplyCache replycache;
// ===================Communication with kcm====================================

//--------------------------1、SendRequest---------------------------------------------
bool isQueueEmpty(){
    // 1=empty,0=not empty
    if (cmdqueue.head == cmdqueue.tail){
        tlogd("cmdqueue is empty,nothing should be sent.\n");
        return true;
    }
    return false;
}

void dequeue(CmdNode *cmdnode){
    cmdnode = &cmdqueue.queue[cmdqueue.head];
    cmdqueue.head = (cmdqueue.head + 1) % MAX_QUEUE_SIZE;
};

//generate Data encryption symmetric key kcm-pub key
TEE_Result generateCmdDataKey(TEE_ObjectHandle key_obj){
    TEE_Result ret;
    TEE_Attribute attr = {0};

    ret = TEE_AllocateTransientObject(TEE_TYPE_AES, AES_KEY_SIZE, &key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate aes transient object, ret 0x%x\n", ret);
        return ret;
    }

    ret = TEE_GenerateKey(key_obj, AES_KEY_SIZE, &attr, 1);
    if (ret != TEE_SUCCESS) {
        tloge("fail to generate aes key, ret 0x%x\n", ret);
        return ret;
    }
}

TEE_Result encryptCmd(uint8_t *jsoncmd, uint32_t jsoncmd_size, TEE_ObjectHandle *key_obj, uint8_t *nonce_buff, uint8_t *tag_buff) {
    TEE_Result ret;
    TEE_OperationHandle oper_enc = NULL;
    uint8_t *encrypted_buffer = NULL;
    size_t encrypted_size = jsoncmd_size;
    size_t nonce_size = NONCE_SIZE;
    size_t tag_size = TAG_SIZE;

    ret = TEE_AllocateOperation(&oper_enc, TEE_ALG_AES_GCM, TEE_MODE_ENCRYPT, AES_KEY_SIZE);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate aes encrypt operation, ret 0x%x\n", ret);
        return ret;
    }

    ret = TEE_SetOperationKey(oper_enc, key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to set rsa encrypt key, ret 0x%x\n", ret);
        TEE_FreeOperation(oper_enc);
        return ret;
    }
    TEE_GenerateRandom(nonce_buff, nonce_size);
    TEE_AEInit(oper_enc, nonce_buff, nonce_size, tag_size, 0, 0);
    encrypted_buffer = TEE_Malloc(encrypted_size*sizeof(uint8_t)+1, 0);

    ret = TEE_AEEncryptFinal(oper_enc, jsoncmd, jsoncmd_size, encrypted_buffer, encrypted_size, tag_buff, &tag_size);
    if (ret != TEE_SUCCESS) {
        tloge("fail to final aes encrypt, ret 0x%x\n", ret);
        TEE_FreeOperation(oper_enc);
        return ret;
    }
    memcpy_s(jsoncmd, jsoncmd_size, encrypted_buffer, encrypted_size);
    TEE_Free(encrypted_buffer);
    return TEE_SUCCESS;
}

TEE_Result encryptKey(TEE_ObjectHandle key_obj, uint8_t *encrypted_key){
    TEE_Result ret;
    TEE_ObjectHandle rsa_key_obj = NULL;
    uint8_t modulus[RSA_KEY_SIZE] = {0};
    uint32_t rsa_size = RSA_KEY_SIZE;
    uint32_t enc_key_len = AES_KEY_SIZE;
    TEE_Attribute attrs[2];

    ret = TEE_AllocateTransientObject(TEE_TYPE_RSA_PUBLIC_KEY, RSA_KEY_SIZE, &rsa_key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate rsa transient object, ret 0x%x\n", ret);
        return ret;
    }
    ret = restoreKeyandCert("sec_storage_data/kcmpub.txt", modulus, &rsa_size);
    if (ret != TEE_SUCCESS) {
        tloge("fail to restore kcm public key, ret 0x%x\n", ret);
        return ret;
    }

    TEE_InitRefAttribute(&attrs[0], TEE_ATTR_RSA_MODULUS, modulus, rsa_size);
    TEE_InitRefAttribute(&attrs[1], TEE_ATTR_RSA_PUBLIC_EXPONENT, exponent, 3);
    ret = TEE_PopulateTransientObject(rsa_key_obj, attrs, sizeof(attrs)/sizeof(TEE_Attribute));
    
    TEE_OperationHandle oper_key_enc = NULL;

    ret = TEE_AllocateOperation(&oper_key_enc, TEE_ALG_RSAES_PKCS1_V1_5, TEE_MODE_ENCRYPT, RSA_KEY_SIZE);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate rsa encrypt operation, ret 0x%x\n", ret);
        return ret;
    }
    ret = TEE_SetOperationKey(oper_key_enc, rsa_key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to set rsa encrypt key, ret 0x%x\n", ret);
        TEE_FreeOperation(oper_key_enc);
        return ret;
    }
    ret = TEE_AsymmetricEncrypt(oper_key_enc, NULL, 0, key_obj->Attribute->content.ref.buffer, AES_KEY_SIZE, encrypted_key, &enc_key_len);
    if (ret != TEE_SUCCESS)
        tloge("Fail to do rsa encrypt, ret 0x%x\n", ret);

    TEE_FreeOperation(oper_key_enc);
    return TEE_SUCCESS;
}

void hex2char(uint32_t hex, uint8_t *hexchar, int32_t i) {
    for(i--; i >= 0; i--, hex >>= 4) {
        if ((hex & 0xf) <= 9)
            *(hexchar + i) = (hex & 0xf) + '0';
        else
            *(hexchar + i) = (hex & 0xf) + 'A' - 0x0a;
    }
}

void uuid2char(TEE_UUID uuid, uint8_t charuuid[33]) {
    int32_t i = 0;

    hex2char(uuid.timeLow, charuuid, 8);
    hex2char(uuid.timeMid, charuuid + 8, 4);
    hex2char(uuid.timeHiAndVersion, charuuid + 12, 4);
    for(i = 0; i < 2; i++){
        hex2char(uuid.clockSeqAndNode[i], charuuid + 16 + i * 2, 2);
    }
    for(i = 0; i < 6; i++){
        hex2char(uuid.clockSeqAndNode[i+2], charuuid + 20 + i * 2, 2);
    }
}

// transfer struct CmdNode to cJSON format
void cmdNode2cjson(CmdNode cmdnode, cJSON *cj) {
    uint8_t taid[33] = {0};
    uint8_t keyid[33] = {0};
    uint8_t masterkey[33] = {0};
    uint8_t ktauuid[33] = {0};
    // translate cmd
    cJSON_AddNumberToObject(cj, "cmd", cmdnode.cmd);
    // translate taid
    uuid2char(cmdnode.taId, taid[33]);
    cJSON_AddStringToObject(cj, "taid", taid[33]);
    // translate keyid
    uuid2char(cmdnode.keyId, keyid[33]);
    cJSON_AddStringToObject(cj, "keyid", keyid[33]);
    // translate masterkey
    uuid2char(cmdnode.masterkey, masterkey[33]);
    cJSON_AddStringToObject(cj, "masterkey", masterkey[33]);
    // translate kta uuid
    uuid2char(ktaUuid, ktauuid[33]);
    cJSON_AddStringToObject(cj, "ktauuid", ktauuid[33]);
}

TEE_Result generateFinalRequest(CmdNode cmdnode, uint8_t *finalrequest){
    TEE_Result ret;
    TEE_ObjectHandle data_key = NULL;
    cJSON cmdjsonnode = {0};
    cJSON finalcmdjsonnode = {0};
    uint8_t nonce_buff[NONCE_SIZE+1] = {0};
    uint8_t tag_buff[TAG_SIZE+1] = {0};
    uint8_t encrypted_key[RSA_KEY_SIZE+1] = {0};
    /* get request data from cmdqueue,and generate final request*/
    ret = generateCmdDataKey(data_key);
    if (ret != TEE_SUCCESS) {
        tloge("fail to generate key, ret 0x%x\n", ret);
        return ret;
    }
    cmdNode2cjson(cmdnode, &cmdjsonnode);//cmdqueue.queue[cmdqueue.head]
    uint8_t *charRequest = cJSON_PrintUnformatted(&cmdjsonnode);
    ret = encryptCmd(charRequest, strlen(charRequest), data_key, nonce_buff, tag_buff);
    if (ret != TEE_SUCCESS) {
        tloge("encrypt cmd failed");
        return ret;
    }
    ret = encryptKey(data_key, encrypted_key);
    if (ret != TEE_SUCCESS) {
        tloge("encrypt key failed");
        return ret;
    }
    // translate key
    cJSON_AddStringToObject(&finalcmdjsonnode, "key", encrypted_key);
    // translate key_size
    cJSON_AddNumberToObject(&finalcmdjsonnode, "key_size", strlen(encrypted_key));
    cJSON_AddStringToObject(&finalcmdjsonnode, "nonce", nonce_buff);
    cJSON_AddStringToObject(&finalcmdjsonnode, "tag", tag_buff);
    // translate cmddata
    cJSON_AddStringToObject(&finalcmdjsonnode, "cmddata", charRequest);
    // translate data_size
    cJSON_AddNumberToObject(&finalcmdjsonnode, "data_size", strlen(charRequest));
    finalrequest = cJSON_PrintUnformatted(&finalcmdjsonnode);
    return TEE_SUCCESS;
}

TEE_Result SendRequest(uint32_t param_type, TEE_Param params[PARAM_COUNT]) {
    TEE_Result ret;
    CmdNode curNode;
    uint8_t *finalrequest = NULL;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_OUTPUT,
        TEE_PARAM_TYPE_VALUE_OUTPUT,
        TEE_PARAM_TYPE_NONE,
        TEE_PARAM_TYPE_NONE)) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    if (params[0].memref.buffer == NULL ||
        params[0].memref.size == 0) {
        tloge("Bad expected parameter value");
        return TEE_ERROR_BAD_PARAMETERS;
    }
    //Judge whether cmd queue is empty 
    if (isQueueEmpty()){
        params[1].value.a = 0;
        tlogd("cmd queue is empty");
        return TEE_SUCCESS;
    }
    //if is not empty, dequeue a request from the queue
    params[1].value.a = (cmdqueue.tail + MAX_QUEUE_SIZE - cmdqueue.head) % MAX_QUEUE_SIZE;
    dequeue(&curNode);

    //generate Request Return value for ka
    ret = generateFinalRequest(curNode, finalrequest);
    if (ret != TEE_SUCCESS) {
        tloge("fail to generate final request");
        return TEE_ERROR_OVERFLOW;
    }
    errno_t err = memcpy_s(params[0].memref.buffer, sizeof(finalrequest), &finalrequest, sizeof(finalrequest));
    if(err != 0) {
        tloge("buffer is too short");
        return TEE_ERROR_SHORT_BUFFER;
    }
    return TEE_SUCCESS;
}

//--------------------------2、GetResponse---------------------------------------------

TEE_Result decryptkey(uint8_t encrypted_key[RSA_KEY_SIZE+1], uint8_t decrypted_key[RSA_KEY_SIZE]){
    TEE_Result ret;
    TEE_ObjectHandle rsa_key_obj = NULL;
    uint8_t modulus[RSA_KEY_SIZE] = {0};
    uint8_t privateExponent[RSA_KEY_SIZE] = {0};
    uint32_t rsa_size = RSA_KEY_SIZE;
    TEE_Attribute attrs[3];
    
    ret = restoreKTAPriv("sec_storage_data/ktakey.txt", modulus[RSA_KEY_SIZE], privateExponent[RSA_KEY_SIZE]);
    if (ret != TEE_SUCCESS) {
        tloge("restore kta private key failed");
        return TEE_ERROR_STORAGE_EIO;
    }

    ret = TEE_AllocateTransientObject(TEE_TYPE_RSA_PUBLIC_KEY, RSA_KEY_SIZE, &rsa_key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate rsa transient object, ret 0x%x\n", ret);
        return ret;
    }
    TEE_InitRefAttribute(&attrs[0], TEE_ATTR_RSA_MODULUS, modulus, rsa_size);
    TEE_InitRefAttribute(&attrs[1], TEE_ATTR_RSA_PUBLIC_EXPONENT, exponent, 3);
    TEE_InitRefAttribute(&attrs[2], TEE_ATTR_RSA_PRIVATE_EXPONENT, privateExponent, rsa_size);

    TEE_OperationHandle oper_key_enc = NULL;

    ret = TEE_AllocateOperation(&oper_key_enc, TEE_ALG_RSAES_PKCS1_V1_5, TEE_MODE_DECRYPT, RSA_KEY_SIZE);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate rsa encrypt operation, ret 0x%x\n", ret);
        return ret;
    }
    ret = TEE_SetOperationKey(oper_key_enc, rsa_key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to set rsa encrypt key, ret 0x%x\n", ret);
        TEE_FreeOperation(oper_key_enc);
        return ret;
    }
    ret = TEE_AsymmetricDecrypt(oper_key_enc, NULL, 0, encrypted_key, RSA_KEY_SIZE, decrypted_key, &rsa_size);
    if (ret != TEE_SUCCESS)
        tloge("Fail to do rsa encrypt, ret 0x%x\n", ret);

    TEE_FreeOperation(oper_key_enc);
    return TEE_SUCCESS;
}

TEE_Result decryptcmd(uint8_t *decrypted_key, uint8_t *encrypted_cmd, uint8_t *decrypted_cmd) {
    TEE_Result ret;
    TEE_OperationHandle oper_dec = NULL;
    uint8_t *decrypted_buffer = NULL;
    size_t encrypted_size = strlen(encrypted_cmd);
    size_t decrypted_size = strlen(encrypted_cmd);
    uint8_t nonce_buf[12] = {0};
    size_t nonce_size = NONCE_SIZE;
    uint8_t tag_buf[16] = {0};
    size_t tag_size = TAG_SIZE;
    uint8_t key_buf[AES_KEY_SIZE] = {0};
    TEE_ObjectHandle key_obj = NULL;
    TEE_Attribute attr = {0};

    ret = TEE_AllocateTransientObject(TEE_TYPE_AES, AES_KEY_SIZE, &key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate aes transient object, ret 0x%x\n", ret);
        return ret;
    }

    ret = TEE_GenerateKey(key_obj, AES_KEY_SIZE, &attr, 1);
    if (ret != TEE_SUCCESS) {
        tloge("fail to generate aes key, ret 0x%x\n", ret);
        return ret;
    }

    ret = TEE_AllocateOperation(&oper_dec, TEE_ALG_AES_GCM, TEE_MODE_DECRYPT, AES_KEY_SIZE);
    if (ret != TEE_SUCCESS) {
        tloge("fail to allocate aes encrypt operation, ret 0x%x\n", ret);
        return ret;
    }

    ret = TEE_SetOperationKey(oper_dec, key_obj);
    if (ret != TEE_SUCCESS) {
        tloge("fail to set rsa encrypt key, ret 0x%x\n", ret);
        TEE_FreeOperation(oper_dec);
        return ret;
    }
    memcpy_s(nonce_buf, nonce_size, decrypted_key, nonce_size);
    memcpy_s(key_buf, AES_KEY_SIZE, decrypted_key+12, AES_KEY_SIZE);
    memcpy_s(tag_buf, tag_size, decrypted_key, tag_size);
    TEE_AEInit(oper_dec, nonce_buf, nonce_size, tag_size, 0, 0);

    ret = TEE_AEDecryptFinal(oper_dec, encrypted_cmd, encrypted_size, decrypted_cmd, &decrypted_size, tag_buf, &tag_size);
    if (ret != TEE_SUCCESS) {
        tloge("fail to final aes decrypt, ret 0x%x\n", ret);
        TEE_FreeOperation(oper_dec);
        return ret;
    }
    TEE_FreeOperation(oper_dec);
    return TEE_SUCCESS;
}

void char2uuid(uint8_t charuuid[33], TEE_UUID *uuid) {
    int32_t i = 0;
    uint8_t *stop = NULL;
    uint8_t buffer[3] = {0};
    uuid->timeLow = strtoul(charuuid, &stop, 16);
    uuid->timeMid = strtoul(stop, &stop, 16);
    uuid->timeHiAndVersion = strtoul(stop, &stop, 16);
    for(i = 0; i < 2; i++) {
        uuid->clockSeqAndNode[i] = strtoul(charuuid + 16 + i * 2, &stop, 16) >> (8 - i * 8);
    }
    for(i = 0; i < 6; i++) {
        buffer[0] = *(charuuid + 20 + i * 2);
        buffer[1] = *(charuuid + 21 + i * 2);
        uuid->clockSeqAndNode[i + 2] = strtoul(buffer, &stop, 16);
    }
}

void parsejson(uint8_t *decrypted_cmd, uint32_t *cmd, TEE_UUID *taid, TEE_UUID *keyid, uint8_t *key) {
    cJSON *cj = cJSON_Parse(decrypted_cmd);
    uint8_t *chartaid = cJSON_GetObjectItemCaseSensitive(cj, "TAId");
    char2uuid(chartaid, taid);
    uint8_t *charkey = cJSON_GetObjectItemCaseSensitive(cj, "KeyId");
    char2uuid(charkey, keyid);
    uint8_t *charcmd = cJSON_GetObjectItemCaseSensitive(cj, "cmd");
    if(!strncmp(charcmd, "2147483649", 10)) {
        *cmd = 1;
    } else if(!strncmp(charcmd, "2147483650", 10)) {
        *cmd = 2;
    } else if (!strncmp(charcmd, "2147483650", 10)) {
        *cmd = 3;
    }
    key = cJSON_GetObjectItemCaseSensitive(cj, "PlainText");
}

void saveTaKey(TEE_UUID TA_uuid, TEE_UUID keyid, uint8_t *keyvalue) {
    int32_t head = cache.head;
    int32_t nxt = -2;
    if (!CheckUUID(cache.ta[head].id, TA_uuid)) {
        int32_t cur = head;
        nxt = cache.ta[cur].next;
        while (nxt != -1) {
            if (CheckUUID(cache.ta[nxt].id, TA_uuid)) {
                cache.ta[cur].next = cache.ta[nxt].next;
                cache.ta[nxt].next = head;
                cache.head = nxt;
                break;
            }
            cur = nxt;
            nxt = cache.ta[nxt].next;
        }
    }
    TaInfo ta = cache.ta[head];
    int32_t thead = ta.head;
    int32_t cur = thead;
    int32_t nxt = ta.key[cur].next;
    while (nxt != -1) {
        cur = nxt;
        nxt = ta.key[nxt].next;
    }
    for(int32_t i=0; i<MAX_KEY_NUM; i++) {
        if(ta.key[i].next == -1 || i != cur) {
            cache.ta[head].head = i;
            cache.ta[head].key[i].id = keyid;
            memcpy_s(cache.ta[thead].key[i].value, strlen(keyvalue), keyvalue, strlen(keyvalue));
            cache.ta[head].key[i].next = thead;
        }
    }
}

void saveGenReplyCache(TEE_UUID TA_uuid, TEE_UUID keyid, uint8_t *keyvalue) {
    int32_t head = replycache.head;
    int32_t cur = head;
    if(head == -1) {
        goto save;
    }
    int32_t nxt = replycache.list[cur].next;
    while (nxt != -1) {
        cur = nxt;
        nxt = replycache.list[nxt].next;
    }
save:
    for(int32_t i=0; i<MAX_KEY_NUM; i++) {
        if(replycache.list[i].next == -1 || i != cur) {
            replycache.list[i].tag = 1;
            replycache.list[i].keyId = keyid;
            memcpy_s(replycache.list[i].keyvalue, strlen(keyvalue), keyvalue, strlen(keyvalue));
            replycache.list[i].next = -1;
            if(cur != -1){
                replycache.list[cur].next = i;
                replycache.tail = i;
            } else {
                replycache.head = i;
                replycache.tail = i;
            }
            break;
        }
    }
}

void saveDelReplyCache(TEE_UUID TA_uuid, TEE_UUID keyid) {
    int32_t head = replycache.head;
    int32_t cur = head;
    if(head == -1) {
        goto save;
    }
    int32_t nxt = replycache.list[cur].next;
    while (nxt != -1) {
        cur = nxt;
        nxt = replycache.list[nxt].next;
    }
save:
    for(int32_t i=0; i<MAX_KEY_NUM; i++) {
        if(replycache.list[i].next == -1 || i != cur) {
            replycache.list[i].tag = 2;
            replycache.list[i].keyId = keyid;
            replycache.list[i].flag = 1;
            replycache.list[i].next = -1;
            if(cur != -1) {
                replycache.list[cur].next = i;
                replycache.tail = i;
            } else {
                replycache.head = i;
                replycache.tail = i;
            }
        break;
        }
    }
}

TEE_Result GetResponse(uint32_t param_type, TEE_Param params[PARAM_COUNT]) {
    TEE_Result ret;
    cJSON *keybuffer = NULL;
    cJSON *cmdbuffer = NULL;
    uint8_t encrypted_key[RSA_KEY_SIZE+1] = {0};
    uint8_t decrypted_key[RSA_KEY_SIZE] = {0};
    uint8_t *encrypted_cmd = NULL;
    uint8_t *decrypted_cmd = NULL;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_INPUT,
        TEE_PARAM_TYPE_VALUE_OUTPUT,
        TEE_PARAM_TYPE_NONE,
        TEE_PARAM_TYPE_NONE)) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    if (params[0].memref.buffer == NULL || params[0].memref.size == 0) {
        tloge("Bad expected parameter value");
        return TEE_ERROR_BAD_PARAMETERS;
    }
    uint8_t *response = TEE_Malloc(params[0].memref.size*sizeof(uint8_t)+1, 0);
    memcpy_s(response, params[0].memref.size*sizeof(uint8_t)+1, params[0].memref.buffer, params[0].memref.size);
    cJSON *cj = cJSON_Parse(response);
    keybuffer = cJSON_GetObjectItemCaseSensitive(cj, "Key");
    memcpy_s(encrypted_key, RSA_KEY_SIZE+1, keybuffer->valuestring, RSA_KEY_SIZE+1);
    ret = decryptkey(encrypted_key, decrypted_key);
    if (ret != TEE_SUCCESS) {
        tloge("decrypt key failed");
        return ret;
    }
    cmdbuffer = cJSON_GetObjectItemCaseSensitive(cj, "retKeyInfo");
    memcpy_s(encrypted_cmd, strlen(keybuffer->valuestring)*sizeof(uint8_t),
            keybuffer->valuestring, strlen(keybuffer->valuestring)*sizeof(uint8_t));
    ret = decryptcmd(decrypted_key, encrypted_cmd, decrypted_cmd);
    if (ret != TEE_SUCCESS) {
        tloge("decrypt cmd failed");
        return ret;
    }
    TEE_UUID taid = {0}, keyid = {0};
    uint8_t *key = NULL;
    uint32_t cmd = 0;
    parsejson(decrypted_cmd, &cmd, &taid, &keyid, key);
    switch(cmd) {
    case 1:
        saveTaKey(taid, keyid, key);
        saveGenReplyCache(taid, keyid, key);
    case 2:
        saveTaKey(taid,keyid,key);
    case 3:
        saveDelReplyCache(taid, keyid);
    }

}

// ===================Communication with kcm from ta====================================

bool generateKcmRequest(TEE_Param params[PARAM_COUNT]){
    /* when kta can't complete ta-operation in local kta,
    generate a request and insert it in cmdqueue*/
    // 若队列已满，则无法添加新命令
    if (cmdqueue.head == cmdqueue.tail + 1) {
        tloge("cmd queue is already full");
        return false;
    }
    CmdNode *n = params[0].memref.buffer;
    cmdqueue.queue[cmdqueue.tail] = *n;
    cmdqueue.tail = (cmdqueue.tail + 1) % MAX_TA_NUM;
    return true;
}


// ===================Communication with ta=============================================

//---------------------------InitTAKey--------------------------------------------------

TEE_Result GenerateTAKey(uint32_t param_type, TEE_Param params[PARAM_COUNT]) {
    //TEE_UUID TA_uuid, TEE_UUID masterkey, char *account, char *password
    //todo: new a key for ta
    TEE_Result ret;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_INPUT,  
        TEE_PARAM_TYPE_NONE,
        TEE_PARAM_TYPE_VALUE_OUTPUT,  
        TEE_PARAM_TYPE_VALUE_OUTPUT )) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    //params[0].memref.buffer内为输入的cmd结构体
    //params[2]值固定为1
    bool res = generateKcmRequest(params); //生成请求成功或失败的结果存放到params[3]的值中
    if (res) {
        params[3].value.b = 1;
        return TEE_SUCCESS;
    }
    params[3].value.b = 0;
    return TEE_ERROR_OVERFLOW;
}
//---------------------------SearchTAKey------------------------------------------------

void flushCache(TEE_UUID taid, TEE_UUID keyid) {
    /*
    flush the cache according to the LRU algorithm
    support two types of element refresh:
    1.ta sequence;
    2.key sequence;
    */
    int32_t head = cache.head;
    if (!CheckUUID(cache.ta[head].id, taid)) {
        int32_t cur = head;
        int32_t nxt = cache.ta[cur].next;
        while (nxt != -1) {
            if (CheckUUID(cache.ta[nxt].id, taid)) {
                cache.ta[cur].next = cache.ta[nxt].next;
                cache.ta[nxt].next = head;
                cache.head = nxt;
                break;
            }
            cur = nxt;
            nxt = cache.ta[nxt].next;
        }
    }
    TaInfo ta = cache.ta[head];
    head = ta.head;
    if (!CheckUUID(ta.key[head].id, keyid)) {
        int32_t cur = head;
        int32_t nxt = ta.key[cur].next;
        while (nxt != -1) {
            if (CheckUUID(ta.key[nxt].id, keyid)) {
                ta.key[cur].next = ta.key[nxt].next;
                ta.key[nxt].next = head;
                ta.head = nxt;
                break;
            }
            cur = nxt;
            nxt = ta.key[nxt].next;
        }
    }
}

TEE_Result SearchTAKey(uint32_t param_type, TEE_Param params[PARAM_COUNT]) {
    //todo: search a certain ta key, if not exist, call generateKcmRequest(）to generate SearchTAKey request
    TEE_Result ret;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_INPUT,  
        TEE_PARAM_TYPE_MEMREF_OUTPUT,
        TEE_PARAM_TYPE_VALUE_OUTPUT,  
        TEE_PARAM_TYPE_VALUE_OUTPUT )) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    //params[0].memref.buffer内为输入的cmd结构体
    CmdNode *n = params[0].memref.buffer;
    int32_t cur = cache.head;
    while (cur != -1) {
        if (CheckUUID(cache.ta[cur].id, n->taId)) {
            TaInfo ta = cache.ta[cur];
            int32_t idx = ta.head;
            while (idx != -1) {
                if (CheckUUID(ta.key[idx].id, n->keyId)) {
                    params[1].memref.size = sizeof(ta.key[idx].value);
                    params[1].memref.buffer = ta.key[idx].value;
                    params[2].value.a = 0;
                    // 更新cache
                    flushCache(n->taId, n->keyId);
                    return TEE_SUCCESS;
                }
                idx = ta.key[idx].next;
            }
        }
        cur = cache.ta[cur].next;
    }
    params[2].value.a = 1;
    bool res = generateKcmRequest(params);
    if (res) {
        params[3].value.b = 1;
        return TEE_SUCCESS;
    }
    params[3].value.b = 0;
    return TEE_ERROR_OVERFLOW;
}

//----------------------------DestoryKey------------------------------------------------

TEE_Result DestoryKey(uint32_t param_type, TEE_Param params[PARAM_COUNT]) {
    //todo: delete a certain key by calling DeleteTAKey(), then generate a delete key request in TaCache
    TEE_Result ret;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_INPUT,  
        TEE_PARAM_TYPE_MEMREF_OUTPUT,
        TEE_PARAM_TYPE_VALUE_OUTPUT,  
        TEE_PARAM_TYPE_VALUE_OUTPUT)) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    //params[0].memref.buffer内为输入的cmd结构体

    TaInfo regTa;
    int32_t taIndex;
    int32_t targetTaIndex;
    int32_t keyIndex;
    int32_t targetKeyIndex;
    CmdNode *n = params[0].memref.buffer;

    //先对TA进行本地证明，证明通过之后根据UUID和密钥ID查询密钥，然后验证TA的账号密码，验证通过后删除指定密钥，最后向KCM发送删除指定密钥的请求

    //TODO:local verification of TA
    //kta通过ka到ras中获取指定ta基准值，在kta中调用本地证明接口获取ta度量报告，然后在kta中进行验证
    //暂定:比较CmdNode和Cache的UUID和账号密码
    //TODO:verify account and password of TA
    if(!verifyTApasswd(n->taId, n->account, n->password)){                  // def of verifyTApasswd() is in ta_authentiate.c
        tloge("verify Ta password failed");
        return TEE_ERROR_BAD_PARAMETERS;
    }
    taIndex = cache.head;
    while (taIndex != END_NULL && !CheckUUID(n->taId,cache.ta[taIndex].id)) // def of CheckUUID() is in ta_authentiate.c
    {
        //loop
        taIndex = cache.ta[taIndex].next; //move to next one
    }
    targetTaIndex = taIndex;

    //TODO:search key on the basis of UUID and key id   (in Cache.TaInfo)
    regTa = cache.ta[targetTaIndex];
    keyIndex = regTa.head;
    while (keyIndex != END_NULL && !CheckUUID(n->keyId ,regTa.key[keyIndex].id))
    {
        //loop
        keyIndex = regTa.key[keyIndex].next; //move to next one
    }
    if(keyIndex == END_NULL){
        tloge("target key not found");
        return TEE_ERROR_BAD_PARAMETERS;
    }
    targetKeyIndex = keyIndex;

    //TODO:delete certain key                           (in Cache.TaInfo)
    if(regTa.head == targetKeyIndex){
        regTa.head = regTa.key[targetKeyIndex].next;
    }
    else{
        keyIndex = regTa.head;
        while (regTa.key[keyIndex].next != targetKeyIndex){
            //loop
            keyIndex = regTa.key[keyIndex].next; //move to next one
        }
        int32_t nextIndex = regTa.key[targetKeyIndex].next;
        regTa.key[keyIndex].next = nextIndex;
        if(nextIndex == END_NULL){
            regTa.tail = keyIndex;
        }
    }
    cache.ta[targetTaIndex] = regTa;

    //TODO:send request of delete key to KCMS
    bool res = generateKcmRequest(params); //生成请求成功或失败的结果存放到params[3]的值中
    if (res) {
        params[3].value.b = 1;
        return TEE_SUCCESS;
    }
    params[3].value.b = 0;
    return TEE_ERROR_OVERFLOW;
}

//----------------------------GetKcmReply------------------------------------------------

TEE_Result GetKcmReply(uint32_t param_type, TEE_Param params[PARAM_COUNT]){
    TEE_Result ret;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_INPUT,  
        TEE_PARAM_TYPE_MEMREF_OUTPUT,
        TEE_PARAM_TYPE_NONE,  
        TEE_PARAM_TYPE_NONE)) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    CmdNode *n = params[0].memref.buffer;
    if (!verifyTApasswd(n->taId, n->account, n->password)) {
        params[1].value.b = 0;
        return TEE_ERROR_ACCESS_DENIED;
    }
    //params[0].memref.buffer内为输入的cmd结构体
    if (replycache.head == -1 && replycache.tail == -1) {
        tloge("get kcm reply error: reply cache is empty\n");
        return TEE_ERROR_ITEM_NOT_FOUND;
    }
    int32_t cur = replycache.head;
    int32_t pre = -2;
    while (cur != -1) {
        if (CheckUUID(replycache.list[cur].taId, n->taId)) {
            params[1].memref.size = sizeof(ReplyNode);
            params[1].memref.buffer = &replycache.list[cur];
            if (pre == -2) {
                replycache.head = replycache.list[cur].next;
            } else {
                replycache.list[pre].next = replycache.list[cur].next;
            }
            return TEE_SUCCESS;
        }
        pre = cur;
        cur = replycache.list[cur].next;
    }
    tloge("get kcm reply error: reply to ta is not found");
    return TEE_ERROR_ITEM_NOT_FOUND;
}

//----------------------------ClearCache------------------------------------------------

TEE_Result ClearCache(uint32_t param_type, TEE_Param params[PARAM_COUNT]) {
    //todo: clear all ta cache
    TEE_Result ret;
    if (!check_param_type(param_type,
        TEE_PARAM_TYPE_MEMREF_INPUT,  
        TEE_PARAM_TYPE_VALUE_OUTPUT,
        TEE_PARAM_TYPE_NONE,  
        TEE_PARAM_TYPE_NONE)) {
        tloge("Bad expected parameter types, 0x%x.\n", param_type);
        return TEE_ERROR_BAD_PARAMETERS;
    }
    //params[0].memref.buffer内为输入的cmd结构体
    CmdNode *n = params[0].memref.buffer;
    // 验证帐号密码
    if (!verifyTApasswd(n->taId, n->account, n->password)) {
        params[1].value.b = 0;
        return TEE_ERROR_ACCESS_DENIED;
    }

    // cache仅1个元素且命中
    if (CheckUUID(cache.ta[cache.head].id, n->taId) && cache.head == cache.tail) {
        cache.head = END_NULL;
        cache.tail = END_NULL;
        tloge("clear ta cache succeeded.\n");
        params[1].value.b = 1;
        return TEE_SUCCESS;
    }

    // cache仅1个元素且未命中
    if (!CheckUUID(cache.ta[cache.head].id, n->taId) && cache.head == cache.tail) {
        tloge("ta cache not fount.\n");
        params[1].value.b = 0;
        return TEE_ERROR_ITEM_NOT_FOUND;
    }

    // cache有2个或以上元素
    int32_t cur = cache.head;
    if (CheckUUID(cache.ta[cur].id, n->taId)) {
        cache.head = cache.ta[cur].next;
        tloge("clear ta cache succeeded.\n");
        params[1].value.b = 1;
        return TEE_SUCCESS;
    }
    int32_t nxt = cache.ta[cur].next;
    while (nxt != END_NULL) {
        TEE_UUID tmp = cache.ta[nxt].id;
        if (CheckUUID(tmp, n->taId)) {
            cache.ta[cur].next = cache.ta[nxt].next;
            if (nxt == cache.tail) {
                cache.tail = cur;
            }
            tloge("clear ta cache succeeded.\n");
            params[1].value.b = 1;
            return TEE_SUCCESS;
        }
        cur = nxt;
        nxt = cache.ta[nxt].next;
    }
    tloge("ta cache not found.\n");
    params[1].value.b = 0;
    return TEE_ERROR_ITEM_NOT_FOUND;
}