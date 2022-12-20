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
Description: initialize module in kta.
	1. 2022-11-04	leezhenxiang
		define the structures.
    2. 2022-11-18   waterh2o
        redefine some interface
    3. 2022-11-25   waterh2o
        Function implementation
*/

#include <kta_common.h>
#include <tee_mem_mgmt_api.h>
#include <tee_trusted_storage_api.h>
#include <tee_crypto_api.h>
#include <tee_crypto_api.h>
#include <string.h>
#include <securec.h>
#include <cJSON.h>

extern Cache cache;
extern CmdQueue cmdqueue;
extern ReplyCache replycache;

TEE_Result saveKeyandCert(char *name, uint8_t *value, size_t size) {
    uint32_t storageID = TEE_OBJECT_STORAGE_PRIVATE;
    uint32_t w_flags = TEE_DATA_FLAG_ACCESS_WRITE;
    void *create_objectID = name;
    TEE_ObjectHandle persistent_data = NULL;
    TEE_Result ret;
    uint8_t *write_buffer = value;
    ret = TEE_CreatePersistentObject(storageID, create_objectID, strlen(create_objectID), w_flags, TEE_HANDLE_NULL, NULL, 0, (&persistent_data));
    if (ret != TEE_SUCCESS) {
        tloge("Failed to create file: ret = 0x%x\n", ret);
        return ret;
    }

    ret = TEE_WriteObjectData(persistent_data, write_buffer, size);
    if (ret != TEE_SUCCESS) {
        tloge("Failed to write file: ret = 0x%x\n", ret);
        TEE_CloseObject(persistent_data);
        return ret;
    }
    TEE_CloseObject(persistent_data);
    return TEE_SUCCESS;
}

TEE_Result saveKTAPriv(char *name, ktaprivkey *value) {
    uint32_t storageID = TEE_OBJECT_STORAGE_PRIVATE;
    uint32_t w_flags = TEE_DATA_FLAG_ACCESS_WRITE;
    void *create_objectID = name;
    TEE_ObjectHandle persistent_data = NULL;
    TEE_Result ret;
    cJSON *kta_priv_json = cJSON_CreateObject();
    uint8_t *kta_priv = NULL;
    cJSON_AddStringToObject(kta_priv_json, "modulus", (char*)value->modulus);
    cJSON_AddStringToObject(kta_priv_json, "privateExponent", (char*)value->privateExponent);
    kta_priv = (uint8_t*)cJSON_PrintUnformatted(kta_priv_json);
    ret = TEE_CreatePersistentObject(storageID, create_objectID, strlen(create_objectID), w_flags, TEE_HANDLE_NULL, NULL, 0, (&persistent_data));
    if (ret != TEE_SUCCESS) {
        tloge("Failed to create file: ret = 0x%x\n", ret);
        return ret;
    }
    ret = TEE_WriteObjectData(persistent_data, kta_priv, strlen((char*)kta_priv));
    if (ret != TEE_SUCCESS) {
        tloge("Failed to write file: ret = 0x%x\n", ret);
        TEE_CloseObject(persistent_data);
        return ret;
    }
    TEE_CloseObject(persistent_data);
    return TEE_SUCCESS;
}

TEE_Result restoreKeyandCert(char *name, uint8_t *buffer, size_t *buf_len) {
    TEE_Result ret;
    uint32_t storageID = TEE_OBJECT_STORAGE_PRIVATE;
    uint32_t r_flags = TEE_DATA_FLAG_ACCESS_READ;
    void *create_objectID = name;
    TEE_ObjectHandle persistent_data = NULL;
    uint32_t pos = 0;
    uint32_t len = 0;
    char *read_buffer = NULL;
    uint32_t count = 0;
    ret = TEE_OpenPersistentObject(storageID, create_objectID, strlen(create_objectID),r_flags, (&persistent_data));
    if (ret != TEE_SUCCESS) {
        tloge("Failed to open file:ret = 0x%x\n", ret);
        return ret;
    }

    ret = TEE_InfoObjectData(persistent_data, &pos, &len);
    if (ret != TEE_SUCCESS) {
        tloge("Failed to open file:ret = 0x%x\n", ret);
        TEE_CloseObject(persistent_data);
        return ret;
    }

    read_buffer = TEE_Malloc(len + 1, 0);
    if (read_buffer == NULL) {
        tloge("Failed to open file:ret = 0x%x\n", ret);
        TEE_CloseObject(persistent_data);
        return ret;
    }

    /* 读取已存入的数据 */
    ret = TEE_ReadObjectData(persistent_data, read_buffer, len, &count);
    if (ret != TEE_SUCCESS) {
        TEE_CloseObject(persistent_data);
        TEE_Free(read_buffer);
        return ret;
    }
    *buf_len = len;
    int32_t rc = memmove_s(buffer, len, read_buffer, len);
    if (rc != 0) {
        TEE_CloseObject(persistent_data);
        TEE_Free(read_buffer);
        return TEE_ERROR_SECURITY;
    }
    TEE_CloseObject(persistent_data);
    TEE_Free(read_buffer);
    return TEE_SUCCESS;
}

TEE_Result restoreKTAPriv(char *name, uint8_t modulus[RSA_PUB_SIZE], uint8_t privateExponent[RSA_PUB_SIZE]) {
    TEE_Result ret;
    uint32_t storageID = TEE_OBJECT_STORAGE_PRIVATE;
    uint32_t r_flags = TEE_DATA_FLAG_ACCESS_READ;
    void *create_objectID = name;
    TEE_ObjectHandle persistent_data = NULL;
    uint32_t pos = 0;
    uint32_t len = 0;
    uint8_t *read_buffer = NULL;
    uint32_t count = 0;
    ret = TEE_OpenPersistentObject(storageID, create_objectID, strlen(create_objectID),r_flags, (&persistent_data));
    if (ret != TEE_SUCCESS) {
        tloge("Failed to open file:ret = 0x%x\n", ret);
        return ret;
    }

    ret = TEE_InfoObjectData(persistent_data, &pos, &len);
    if (ret != TEE_SUCCESS) {
        tloge("Failed to open file:ret = 0x%x\n", ret);
        TEE_CloseObject(persistent_data);
        return ret;
    }

    read_buffer = TEE_Malloc(len + 1, 0);
    if (read_buffer == NULL) {
        tloge("Failed to open file:ret = 0x%x\n", ret);
        TEE_CloseObject(persistent_data);
        return ret;
    }

    /* 读取已存入的数据 */
    ret = TEE_ReadObjectData(persistent_data, read_buffer, len, &count);
    if (ret != TEE_SUCCESS) {
        TEE_CloseObject(persistent_data);
        TEE_Free(read_buffer);
        return ret;
    }
    cJSON *kta_priv_json = cJSON_Parse((char*)read_buffer);
    cJSON *jsonmodulus = cJSON_GetObjectItemCaseSensitive(kta_priv_json, "modulus");
    cJSON *jsonprivateExponent = cJSON_GetObjectItemCaseSensitive(kta_priv_json, "privateExponent");
    memcpy_s(modulus, RSA_PUB_SIZE, jsonmodulus->valuestring, RSA_PUB_SIZE);
    memcpy_s(privateExponent, RSA_PUB_SIZE, jsonprivateExponent->valuestring, RSA_PUB_SIZE);
    TEE_CloseObject(persistent_data);
    TEE_Free(read_buffer);
    return TEE_SUCCESS;
}

TEE_Result initStructure(){
    //init cache
    cache.head = -1;
    cache.tail = -1;
    for(int i=0;i<MAX_TA_NUM;i++){
        cache.ta[i].next = -1;
        cache.ta[i].head = -1;
        cache.ta[i].tail = -1;
        for(int j=0;j<MAX_KEY_NUM;j++){
            cache.ta[i].key[j].next = -1;
        }
    }
    //init cmdqueue
    cmdqueue.head = 0;
    cmdqueue.tail = 0;
    //init replycache
    replycache.head = -1;
    replycache.tail = -1;
    for(int i=0;i<MAX_QUEUE_SIZE;i++){
        replycache.list[i].next = -1;
    }
    return TEE_SUCCESS;
}

TEE_Result reset(char *name){
    TEE_Result ret;
    uint32_t storageID = TEE_OBJECT_STORAGE_PRIVATE;
    uint32_t r_flags = TEE_DATA_FLAG_ACCESS_READ;
    TEE_ObjectHandle persistent_data = NULL;
    ret = TEE_OpenPersistentObject(storageID, name, strlen(name),
    r_flags | TEE_DATA_FLAG_ACCESS_WRITE_META, (&persistent_data));
    if (ret != TEE_SUCCESS) {
        tloge("Failed to execute TEE_OpenPersistentObject:ret = %x\n", ret);
        return ret;
    }
    TEE_CloseAndDeletePersistentObject(persistent_data);
    return TEE_SUCCESS;
}

TEE_Result Reset_All(){
    TEE_Result ret;
    ret = reset("sec_storage_data/ktacert.txt");
    if (ret != TEE_SUCCESS) {
        tloge("Failed to reset ktacert\n", ret);
        return ret;
    }
    reset("sec_storage_data/kcmpub.txt");
    if (ret != TEE_SUCCESS) {
        tloge("Failed to reset kcmpub\n", ret);
        return ret;
    }
    reset("sec_storage_data/ktakey.txt");
    if (ret != TEE_SUCCESS) {
        tloge("Failed to reset ktakey\n", ret);
        return ret;
    }
    return TEE_SUCCESS;
}
