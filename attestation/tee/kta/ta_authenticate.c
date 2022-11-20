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
Description: ta authenticating module in kta.
	1. 2022-11-04	leezhenxiang
		define the structures.
*/

#include <tee_defines.h>
#include <kta_common.h>

void addTaState(TEE_UUID TA_uuid, char *taId, char *passWd, Cache *cache) {
    //todo: add a ta's id and passwd in to the table
}

void deleteTaState(TEE_UUID TA_uuid, char *taId, char *passWd, Cache *cache) {
    //todo: delete a ta state from tacache
}

void searchTAState(TEE_UUID TA_uuid, char *taId, char *passWd, Cache *cache) {
    //todo: search a ta state from tacache
}

void attestTA() {
    //attest a ta's trusted station locally by QTA
}
