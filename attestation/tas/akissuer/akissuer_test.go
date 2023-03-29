/*
kunpengsecl licensed under the Mulan PSL v2.
You can use this software according to the terms and conditions of
the Mulan PSL v2. You may obtain a copy of Mulan PSL v2 at:
    http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.

Author: wangli/wanghaijing
Create: 2022-04-01
Description: Implement specific services provided by AS
*/

package akissuer

import (
	"io/ioutil"
	"os"
	"testing"

	"gitee.com/openeuler/kunpengsecl/attestation/tas/config"
)

var (
	bufferDAA = []byte(
`
{
        "signature":    {
                "drk_cert":     "TUlJRWtqQ0NBM3FnQXdJQkFnSVJFV0ltOHVOTE5NU3ptZU9yRkhlM0p6Y3dEUVlKS29aSWh2Y05BUUVMQlFBd1BURUxNQWtHQTFVRUJoTUNRMDR4RHpBTkJnTlZCQW9UQmtoMVlYZGxhVEVkTUJzR0ExVUVBeE1VU0hWaGQyVnBJRWxVSUZCeWIyUjFZM1FnUTBFd0hoY05Nakl3TnpBNE1Ea3dOVFE1V2hjTk16Y3dOekEwTURrd05UUTVXakE2TVFzd0NRWURWUVFHRXdKRFRqRVBNQTBHQTFVRUNoTUdTSFZoZDJWcE1Sb3dHQVlEVlFRREV4RXdNalpRVUZZeE1FdEJNREExTVRNNVZEQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0NnZ0lCQUxrSC82NkNTYU52MXhFVEd0WUtHZkNNbzVSTHlXYUdSVmFUbjR2QkhTUDk0ZXhvbU5DMjZSaHV4bEtVK1dNOEw0TjFSVDkxUnZRbVlUNUw3cG1xTzFlb043ei82YUhzZG5wZ3R4RUNSTzZ5WkxGcUs2UEhXd0twUnFkczRlYVcyZlZHVXF2MWZSRmVTaEtGM29YUVpYbmpLN1owS3M0QTZ6TmFNUnBmV3VNdlYxdHBXOFZVQVk0WG5lQjl0Zzk4d3ZRZUorR3k2cVoyVEFYVEkrekhUWnVMNENtTnUrTDdJTkNaT0VWUnNwVm1HVFhtSkhzbjhYQ0dWYllRMk5RT0FkdW9WUnRWaFZSMjk4SHQ0NUc1ejNTQzZtdEY1dTdPeXYvVlFQbytuS3hmRUNxLytjTUpIK3lsQWFCcXlya1hiYkU2SDRieU1PZUN5REFLSzlvSDdsdXNBenA2aURpN1QrRkFXZ3BRNzloMHJ4UUVBN1FhVGNDYXhiUG5RanNtYzd1eHN3M0YwWVdibE4wVEhFNGpTc1krNHZGNGtZaXJRbWNCbE5MSWJIUkhsb1RDQzJZcGh1ck5vNDZjSnFTR3I5Q1pieHZVbkt1U2tQelV5RnpkdEsyUTA2L0pFTXF1RzNMRC9TeGRvVld4RzgyVVhvYTRXaE54MXBRU1ZVdDREVU9rUzk1RVJsTkFVWnpteGczSjBhOC9PTGl2MHlXUm1EQmF5ZnU1WERJbUJaRGJ5WHVUdUQ1a2Y3UmhETjJXTkcwUVFWY0RPK1UyaEtTSWd6cGthL3lNTDVDUW9rL3pXQnM4MWJteUVnOG5vbXRwLy9KSkJocU5DOHpJcEIwWGhoMVMySHdaNCtDK3k1b2Ywd3Ardnl5SGp4RmZuWVR3K3hVRk83VDJjZTdwQWdNQkFBR2pnWTh3Z1l3d0h3WURWUjBqQkJnd0ZvQVVFb28zN1BselY5ZmtnN3ExYndQczNhNFJUczh3Q3dZRFZSMFBCQVFEQWdQNE1Gd0dDQ3NHQVFVRkJ3RUJCRkF3VGpBb0JnZ3JCZ0VGQlFjd0FvWWNhSFIwY0Rvdkx6RXlOeTR3TGpBdU1TOWpZV2x6YzNWbExtaDBiVEFpQmdnckJnRUZCUWN3QVlZV2FIUjBjRG92THpFeU55NHdMakF1TVRveU1EUTBNekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBWTZaaFdnb1IxbWxVT0JMbGJCRjBGZ0l2eFdxcFJld1lYUk1UT1NXTlF2K24xa1FkSFFOTGkrV2hhM21QRENMeHV5N2xCcEI3QTN5cDdvZDIwS3F6U25UcXFDNHExdkZqVGxpcmlxYVZvcjg4bXhpN2wrV2JIeW5VMTduYVpITG5DcEo1Q3RHdmJLREMvRGJGVHJvbDZ1SnBYTmp4eDNzZDBleUw5azNXYUIxWmgvZEV3cUtFdEw3K2xlSDg5MEM3enRTUGsyUkVjbFlLMTl1UWVPdzFNRzhpcVhBZTIrSlNzMUxuUzI3TFlsOG9xaWhDdkM2cGpTK3RkTThLTGk4cUd3MHBNM0lDc3pjQmlhOW5Od0tna0RIWkp1K1JyamtvZ3BGcjlWV3FJWG41cFMvekloc3Y2MEZaSnRKTXpsNnlPbXVYdVM4U1BtYk9Jd3FoblltQ2hnPT0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
                "drk_sign":     "JutEwU_gMbxLnKVEnpsipSytcZnr0d8cCpAoFaN_mE0dzIJvFY2PjievU-gt2peLRm1sI1tkMy_1j4upIsduJ23SMrGw7azo2R0xtgDlEM3vv030A2PAlFsZ2lT-Q6nsmnkj4ddWm4-MBtj2a4YCOazEpTaT1AeuD1KNQrlo6utVyEog_cvrPB40WszVXdZWwTGkfMDDGuRBAP9eH96q7uxqUrg4K2nUb7sOmFbW5SgmaMoILgdEp3vV_7xfROPJigmX35zAdAlMBVSC6kiER7sm5kSrhaC7b2yJjZQtynfXizq1oTKqSokbmg32DHDflKnb0PjUi9xcyRvm4MglgSgsybCH0m6atdWyaqsjPG7txAMYWgwg0zpZMtl7OhllCfi1GS4ZoHTiXetvMIE1wuPBK-Gtt5SOLndIVazNzkYgExVgBtpcA9bs2txM_J2jgZaayP_jZgYRNcihsWwnz8qsxVotNpDfcg8e4aJmPb32_8ji5NvQ4HR4vcrzTAypstMHiJWEkp4SDO8cHoNLhvwE4TDBl9sMZeZORvvIBjpp3F62T_PEZi2aCsq68RrJnO-mi7XZhr17RvVu_y-vxDOlOeGPlymK9nlUqkmzAb9PX7YcpaKLYrjO7M8fQwrN2RSMH9NZ9flRNlHIPTigLMl2GEV5lLJyxSAbE3z_tOc"
        },
        "payload":      {
                "version":      "TEE.RA.1.0",
                "timestamp":    3713734665132146,
                "scenario":     "sce_as_with_daa",
                "sign_alg":     "PS256",
                "hash_alg":     "HS256",
                "qta_img":      "jG-vRNjR-p-Dl35AXKrdfZsRlIOdYDTVQS_gD18pSXY",
                "qta_mem":      "D0dOvFbcVQL3-6KibZg7Omzy2oTORHMPYNTGYNQmNKY",
                "tcb":  "",
                "ak_pub":       {
                        "kty":  "DAA",
                        "qs":   "QAAAAC3b_yop7eXo1inYHuhC0AiTB2xecZ5RgKelKlQZGC_yu0EI22L6z50_sOeSqmhU5se1evZZNyabm7rWmp-XyC5AAAAAUI2JgC4kidxj8ZkjPNIZVFmo0bas2t4cdTaSKv_LsFYN0plo0RQqf78upMVILYroGJwXVS0V9iFPrkexhTsKPA"
                }
        },
        "handler":      "provisioning-output"
}
`)

	bufferNoDAA = []byte(
`
{
        "signature":    {
                "drk_cert":     "TUlJRWtqQ0NBM3FnQXdJQkFnSVJFV0ltOHVOTE5NU3ptZU9yRkhlM0p6Y3dEUVlKS29aSWh2Y05BUUVMQlFBd1BURUxNQWtHQTFVRUJoTUNRMDR4RHpBTkJnTlZCQW9UQmtoMVlYZGxhVEVkTUJzR0ExVUVBeE1VU0hWaGQyVnBJRWxVSUZCeWIyUjFZM1FnUTBFd0hoY05Nakl3TnpBNE1Ea3dOVFE1V2hjTk16Y3dOekEwTURrd05UUTVXakE2TVFzd0NRWURWUVFHRXdKRFRqRVBNQTBHQTFVRUNoTUdTSFZoZDJWcE1Sb3dHQVlEVlFRREV4RXdNalpRVUZZeE1FdEJNREExTVRNNVZEQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0NnZ0lCQUxrSC82NkNTYU52MXhFVEd0WUtHZkNNbzVSTHlXYUdSVmFUbjR2QkhTUDk0ZXhvbU5DMjZSaHV4bEtVK1dNOEw0TjFSVDkxUnZRbVlUNUw3cG1xTzFlb043ei82YUhzZG5wZ3R4RUNSTzZ5WkxGcUs2UEhXd0twUnFkczRlYVcyZlZHVXF2MWZSRmVTaEtGM29YUVpYbmpLN1owS3M0QTZ6TmFNUnBmV3VNdlYxdHBXOFZVQVk0WG5lQjl0Zzk4d3ZRZUorR3k2cVoyVEFYVEkrekhUWnVMNENtTnUrTDdJTkNaT0VWUnNwVm1HVFhtSkhzbjhYQ0dWYllRMk5RT0FkdW9WUnRWaFZSMjk4SHQ0NUc1ejNTQzZtdEY1dTdPeXYvVlFQbytuS3hmRUNxLytjTUpIK3lsQWFCcXlya1hiYkU2SDRieU1PZUN5REFLSzlvSDdsdXNBenA2aURpN1QrRkFXZ3BRNzloMHJ4UUVBN1FhVGNDYXhiUG5RanNtYzd1eHN3M0YwWVdibE4wVEhFNGpTc1krNHZGNGtZaXJRbWNCbE5MSWJIUkhsb1RDQzJZcGh1ck5vNDZjSnFTR3I5Q1pieHZVbkt1U2tQelV5RnpkdEsyUTA2L0pFTXF1RzNMRC9TeGRvVld4RzgyVVhvYTRXaE54MXBRU1ZVdDREVU9rUzk1RVJsTkFVWnpteGczSjBhOC9PTGl2MHlXUm1EQmF5ZnU1WERJbUJaRGJ5WHVUdUQ1a2Y3UmhETjJXTkcwUVFWY0RPK1UyaEtTSWd6cGthL3lNTDVDUW9rL3pXQnM4MWJteUVnOG5vbXRwLy9KSkJocU5DOHpJcEIwWGhoMVMySHdaNCtDK3k1b2Ywd3Ardnl5SGp4RmZuWVR3K3hVRk83VDJjZTdwQWdNQkFBR2pnWTh3Z1l3d0h3WURWUjBqQkJnd0ZvQVVFb28zN1BselY5ZmtnN3ExYndQczNhNFJUczh3Q3dZRFZSMFBCQVFEQWdQNE1Gd0dDQ3NHQVFVRkJ3RUJCRkF3VGpBb0JnZ3JCZ0VGQlFjd0FvWWNhSFIwY0Rvdkx6RXlOeTR3TGpBdU1TOWpZV2x6YzNWbExtaDBiVEFpQmdnckJnRUZCUWN3QVlZV2FIUjBjRG92THpFeU55NHdMakF1TVRveU1EUTBNekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBWTZaaFdnb1IxbWxVT0JMbGJCRjBGZ0l2eFdxcFJld1lYUk1UT1NXTlF2K24xa1FkSFFOTGkrV2hhM21QRENMeHV5N2xCcEI3QTN5cDdvZDIwS3F6U25UcXFDNHExdkZqVGxpcmlxYVZvcjg4bXhpN2wrV2JIeW5VMTduYVpITG5DcEo1Q3RHdmJLREMvRGJGVHJvbDZ1SnBYTmp4eDNzZDBleUw5azNXYUIxWmgvZEV3cUtFdEw3K2xlSDg5MEM3enRTUGsyUkVjbFlLMTl1UWVPdzFNRzhpcVhBZTIrSlNzMUxuUzI3TFlsOG9xaWhDdkM2cGpTK3RkTThLTGk4cUd3MHBNM0lDc3pjQmlhOW5Od0tna0RIWkp1K1JyamtvZ3BGcjlWV3FJWG41cFMvekloc3Y2MEZaSnRKTXpsNnlPbXVYdVM4U1BtYk9Jd3FoblltQ2hnPT0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
                "drk_sign":     "KxpJcxJD6pTdZRbpcYkR3pj4XlNTmCCKAbAWGh4QnJcaaiah6L08hHjamS5FCG24gj5VulbOphQ_UFricLms9gNsQi2gD9hLT7UG3HhzDrfcMcbzNBSM58bJ0qBt9YYsRlJsm-a6RE3VLmOMtkcj5KkpLPIuXvVwXa_HpR2Kck3nmzYp08f4wA7aGUC_w7AS4ITnCxVk3phAIsGHmoSK7wW96Dn80HiVlOWzUzESwp3c79al18MvVDJFv-2DRUIegfpjPsGvJiU5fcwIqRi5ebn-aCrcS10s3TWVDJLHj1pJCd8Fy07AuIvkoYgI_3mjTC-FbrgERLdctaiVZiiG8wS_w8HJalzq3Hb6BE6TYxbtD6ex6GDt09YQxv51KhiCbgY8qgqNgqXtZgsZCVPUf-2-yrH7czd1w3Rj4goCHtAakLQMma1OgPFr_6ScvZGqO0nuv7Un-keJc938Jq1yN8IAISq3skedCVJLqRzZpQYbOjjb3afj5DEiLkqSwU5aDZ0ecEHwIBjJqmwYjZDuOAsKQ_idOWjGmZ9fNwUm_UTkEEPs5jBFH0pVjNb8ZsdTCF2X3V4DebwE0JniNBd0MnFyb_mQY_NP58b0F94I8P5-e4hM2DiZ0wGDAe-ungKkDLYX8VrFUdgpScZwiAJ9tEq-oJ2Qs8AHKLuzFIwJwwI"
        },
        "payload":      {
                "version":      "TEE.RA.1.0",
                "timestamp":    3710375507375324,
                "scenario":     "sce_as_no_daa",
                "sign_alg":     "PS256",
                "hash_alg":     "HS256",
                "qta_img":      "jG-vRNjR-p-Dl35AXKrdfZsRlIOdYDTVQS_gD18pSXY",
                "qta_mem":      "D0dOvFbcVQL3-6KibZg7Omzy2oTORHMPYNTGYNQmNKY",
                "tcb":  "",
                "ak_pub":       {
                        "kty":  "RSA",
                        "n":    "3FUp1JfTJH9Nb46eVrPdfeuf-lht-K0eRV_XDqVbqf0_45A5EDiPYAwsZGY8gs2zVmS1Kpt1PioRgbvzJ97urkuOGVBfrHUW82gSPXsgbj6rKg2jAg2uXdEd7pObVWURkxsvtHw4QYVx9yru9EWtFaHfqVIpcJdI16nPHAKv_Wal_fm82rCxaw5Es-8xk5a0Yytu8G3n940fJ7UW2mnLBXWKXoThi0Ad6yahxtxac3p5qWIb8Mk6tk1N_kIPdu4HrsQPdAv3NuUOZ0mA6U9EQYh10BOoA7NUUr527jD116XnoN8OdiGYTs5WB_wz3FOa3gbwukctlTA33VasTL-JqGg0nRPukrr8HDLgUV3Xr5QFwLv2loVCtbyPsfpQnOX3iBnU894T75zwcJSy5TAKZB1DuDFLik-qyukJcoW6dVt9b0aY4QO6ajzMJr4p8ZPAaFtUg1N3rFsCCWrQbhjZpJM_19Mr1FN28_s9FmLTSFZUpjsJFd74ZTeaKFJqso1rXOTL2HE1OOejaPSAhyZSo4CC7Js-2VNMYj5aIgU1qrg9V9shfyhCSHMkVvymcJ4az2_HJxmzOUrXkhEqLGH48nWj_bSZ9DF4d3ms4L12GUs9pdorgt1vz9mufC4raGE7haRNIksmbarJbaPrx4ERVAbDhcSavPmgjfgG7E7Leis",
                        "e":    "AQAB"
                }
        },
        "handler":      "provisioning-output"
}
`)
)

const serverConfig = `
tasconfig:
  port: 127.0.0.1:40008
  rest: 127.0.0.1:40009
  akskeycertfile: ./ascert.crt
  aksprivkeyfile: ./aspriv.key
  huaweiitcafile: ./Huawei IT Product CA.pem
  DAA_GRP_KEY_SK_X: 65a9bf91ac8832379ff04dd2c6def16d48a56be244f6e19274e97881a776543c65a9bf91ac8832379ff04dd2c6def16d48a56be244f6e19274e97881a776543c
  DAA_GRP_KEY_SK_Y: 126f74258bb0ceca2ae7522c51825f980549ec1ef24f81d189d17e38f1773b56126f74258bb0ceca2ae7522c51825f980549ec1ef24f81d189d17e38f1773b56
  basevalue: "8c6faf44d8d1fa9f83977e405caadd7d9b1194839d6034d5412fe00f5f294976 0f474ebc56dc5502f7fba2a26d983b3a6cf2da84ce44730f60d4c660d42634a6"
`

const (
	asCertPath = "./ascert.crt"
	asprivPath = "./aspriv.key"
	huaweiPath = "./Huawei IT Product CA.pem"
)

const ascert = `
-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIURBG3rzn2SH3wruHaPKctA2pEqS8wDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMjEyMTcxMzE3MTRaFw0yMzEy
MTcxMzE3MTRaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQC3FM8ErX6+Vcy4sf4DQexeLDRFrqaKOWEEi90MTjE+
aAuIV0Vy8tFgnvaVoznLD6J7LUhaT1xvVNEhsQpSp+BEdZ8kUtxNm91j4FBDE+V6
RxWbGPXcZhfVSmoc+dlwt5aCrAH03GAJOp6NWLBGA6LBN1X5yaT6hJJTq3ioC+K3
h/08Ub801Glnai4BECygUHp1qn7bFl9p2o5yRfL80MJxQKsq1dKvUY/+Dsxzdwcr
GVT68nstBjXbnoTAjyhIZepJvZmGKvoiEG0sQwi3dInEK5Q9QLZBAHPNS+M/3v9g
dW/F1QPYeFEtz7+Np9gUYs9GOkB5L6Cr+gXeFRfRiO1/g9BMzrFhkKYTgnLHpdwu
oWpDmmfcr8TcRgSZCZ9oAC41wcn5b8nuGNUbI8fGRKaDzEsprEID5gOVofsC2KZv
sD+9ATI3MBD95OKXIbE0zMaAx21EIUxe7HVWUTcNnHf9kBpVR8x27xa/AV4Qc89d
1pCYg0cm/4SpkuhzfP1VGdL935QQcVumR85nBNyX4l8c8239zbnIepUn8GzkRqi5
5dwqLzxAY96y2VYNT3pbuRd3qK+psuPlviSZ/NtebFETjzVbl6LRCsKmY2DknfBo
LyIEO7M06H6B15lbZjPEd0vil4R9gROl/7h9k8wP8rj68/BiCE32xvsXRyy5oDt4
bwIDAQABo1MwUTAdBgNVHQ4EFgQUv+5bJBMwfJN9pB0XhB+niXZu6lowHwYDVR0j
BBgwFoAUv+5bJBMwfJN9pB0XhB+niXZu6lowDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEACAIBde61vKD3PBLiYmgFtXErnRIkGm9EY8q/T74xv9vX
b04mry+LUHAKDx/M2wpfcGW2rAGaNGvgGvfhK/vKv9P7gNmIjZOGgSJm+lsKCr39
2NlROMsi08GGWRBQhZNEt5feaH5bcCGWjDHnNTL1Nhe/OOf1i74X7gX3WS1mD0O+
I9/TUznmNg7bZhICRswFHEymSHMxyOsvzG+f1ENUr6XKgXTWD89PNOJ0IzQsXq7V
W96YSM7EvW87AXWyioFi7B9TRHtSxK+/ZJz5joZos8X4/Yamve7OpX3jQnrxxh0W
vkNdJ1fiiYzEciyTHAVUTA1q/ZqewEUgVZYhIAbTCEV1h1PLrL3VHbzCrarWZqX1
+vSDOJoYBtAMfugcsYqgnIdYOSwpQjdan8rXIqhgk2rwAgmIjEWvRAvFhoOxO7Os
PhLgoeJo+JahmUjAfE8L/wW3k/3OJQy2eD+cAtUFWpOdMrmE5FtepUv6voM6N6Fa
Q+d69RUiN0XvVYlG/ZXIeHtPYk6cXxof7J3Tn8ieEKbT7OwuG7vxRodCGSao6o3i
uOdEBUnNkMadj7i265D8de8sOuQ7+pPu4lNBZGmF5XeaBEEynF8ts1rWrwM+lUgK
/bK1vXGMemom0NEsj6zAOd2GUBuhFP9WLYXh4SRTY/5PKVhAdL9oenaIxji99lQ=
-----END CERTIFICATE-----
`

const aspriv = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

const huaweica = `
-----BEGIN CERTIFICATE-----
MIIEsTCCApmgAwIBAgIRdjl5z9FobnagzdStBIQZVIcwDQYJKoZIhvcNAQELBQAw
PDELMAkGA1UEBhMCQ04xDzANBgNVBAoTBkh1YXdlaTEcMBoGA1UEAxMTSHVhd2Vp
IEVxdWlwbWVudCBDQTAeFw0xNjEwMTgwNjUwNTNaFw00MTEwMTIwNjUwNTNaMD0x
CzAJBgNVBAYTAkNOMQ8wDQYDVQQKEwZIdWF3ZWkxHTAbBgNVBAMTFEh1YXdlaSBJ
VCBQcm9kdWN0IENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtKE3
0649koONgSJqzwKXpSxTwiGTGorzcd3paBGH75Zgm5GFv2K2TG3cU6seS6dt7Ig+
/8ntrcieQUttcWxpm2a1IBeohU1OTGFpomQCRqesDnlXXUS4JgZiDvPBzoqGCZkX
YRw37J5KM5TSZzdLcWgxAPjXvKPdLXfxGzhqg8GV1tTboqXoNEqVqOeViBjsjN7i
xIuu1Stauy9E0E5ZnSrwUjHc5QrR9CmWIu9D0ZJJp1M9VgcXy9evPhiHoz9o+KBd
fNwt4e/NymTqaPa+ngS/qZwI7A4tR4RKCMKFHJcsjaXwUb0RuIeCiPO3wPHgXmGL
uiKfyPV8SMLpE/wYaQIDAQABo4GsMIGpMB8GA1UdIwQYMBaAFCr4EFkngDUfp3y6
O58q5Eqqm5LqMEYGA1UdIAQ/MD0wOwYEVR0gADAzMDEGCCsGAQUFBwIBFiVodHRw
Oi8vc3VwcG9ydC5odWF3ZWkuY29tL3N1cHBvcnQvcGtpMA8GA1UdEwQIMAYBAf8C
AQAwDgYDVR0PAQH/BAQDAgEGMB0GA1UdDgQWBBQSijfs+XNX1+SDurVvA+zdrhFO
zzANBgkqhkiG9w0BAQsFAAOCAgEAAg1oBG8YFvDEecVbhkxU95svvlTKlrb4l77u
cnCNhbnSlk8FVc5CpV0Q7SMeBNJhmUOA2xdFsfe0eHx9P3Bjy+difkpID/ow7oBH
q2TXePxydo+AxA0OgAvdgF1RBPTpqDOF1M87eUpJ/DyhiBEE5m+QZ6VqOi2WCEL7
qPGRbwjAFF1SFHTJMcxldwF6Q/QWUPMm8LUzod7gZrgP8FhwhDOtGHY5nEhWdADa
F9xKejqyDCLEyfzsBKT8V4MsdAo6cxyCEmwiQH8sMTLerwyXo2o9w9J7+vRAFr2i
tA7TwGF77Y1uV3aMj7n81UrXxqx0P8qwb467u+3Rj2Cs29PzhxYZxYsuov9YeTrv
GfG9voXz48q8ELf7UOGrhG9e0yfph5UjS0P6ksbYInPXuuvrbrDkQvLBYb9hY78a
pwHn89PhRWE9HQwNnflTZS1gWtn5dQ4uvWAfX19e87AcHzp3vL4J2bCxxPXEE081
3vhqtnU9Rlv/EJAMauZ3DKsMMsYX8i35ENhfto0ZLz1Aln0qtUOZ63h/VxQwGVC0
OCE1U776UUKZosfTmNLld4miJnwsk8AmLaMxWOyRsqzESHa2x1t2sXF8s0/LW5T7
d+j7JrLzey3bncx7wceASUUL3iAzICHYr728fNzXKV6OcZpjGdYdVREpM26sbxLo
77rH32o=
-----END CERTIFICATE-----
`

const (
	configFilePath = "./config.yaml"
)

func createFiles() {
	ioutil.WriteFile(configFilePath, []byte(serverConfig), 0644)
	ioutil.WriteFile(asCertPath, []byte(ascert), 0644)
	ioutil.WriteFile(asprivPath, []byte(aspriv), 0644)
	ioutil.WriteFile(huaweiPath, []byte(huaweica), 0644)
}

func removeFiles() {
	os.Remove(configFilePath)
	os.Remove(asCertPath)
	os.Remove(asprivPath)
	os.Remove(huaweiPath)
}

func TestGenerateDAAAKCert(t *testing.T) {
	createFiles()
	defer removeFiles()

	config.LoadConfigs()

	err := config.InitializeAS()
	if err != nil {
		t.Error(err)
	}
	_, _, err = GenerateDAAAKCert(bufferDAA)
	if err != nil {
		t.Errorf("generate daa scenario ak cert error %v", err)
	}
}

func TestGenerateNoDAAAKCert(t *testing.T) {
	createFiles()
	defer removeFiles()

	config.LoadConfigs()

	err := config.InitializeAS()
	if err != nil {
		t.Error(err)
	}
	_, err = GenerateNoDAAAKCert(bufferNoDAA)
	if err != nil {
		t.Errorf("generate nodaa scenario ak cert error %v", err)
	}
}

func TestGetNoDAAData(t *testing.T) {
	jsonData := []byte(`{
		"signature":	{
			"drk_cert":	"TUlJRWtqQ0NBM3FnQXdJQkFnSVJFV0ltOHVOTE5NU3ptZU9yRkhlM0p6Y3dEUVlKS29aSWh2Y05BUUVMQlFBd1BURUxNQWtHQTFVRUJoTUNRMDR4RHpBTkJnTlZCQW9UQmtoMVlYZGxhVEVkTUJzR0ExVUVBeE1VU0hWaGQyVnBJRWxVSUZCeWIyUjFZM1FnUTBFd0hoY05Nakl3TnpBNE1Ea3dOVFE1V2hjTk16Y3dOekEwTURrd05UUTVXakE2TVFzd0NRWURWUVFHRXdKRFRqRVBNQTBHQTFVRUNoTUdTSFZoZDJWcE1Sb3dHQVlEVlFRREV4RXdNalpRVUZZeE1FdEJNREExTVRNNVZEQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0NnZ0lCQUxrSC82NkNTYU52MXhFVEd0WUtHZkNNbzVSTHlXYUdSVmFUbjR2QkhTUDk0ZXhvbU5DMjZSaHV4bEtVK1dNOEw0TjFSVDkxUnZRbVlUNUw3cG1xTzFlb043ei82YUhzZG5wZ3R4RUNSTzZ5WkxGcUs2UEhXd0twUnFkczRlYVcyZlZHVXF2MWZSRmVTaEtGM29YUVpYbmpLN1owS3M0QTZ6TmFNUnBmV3VNdlYxdHBXOFZVQVk0WG5lQjl0Zzk4d3ZRZUorR3k2cVoyVEFYVEkrekhUWnVMNENtTnUrTDdJTkNaT0VWUnNwVm1HVFhtSkhzbjhYQ0dWYllRMk5RT0FkdW9WUnRWaFZSMjk4SHQ0NUc1ejNTQzZtdEY1dTdPeXYvVlFQbytuS3hmRUNxLytjTUpIK3lsQWFCcXlya1hiYkU2SDRieU1PZUN5REFLSzlvSDdsdXNBenA2aURpN1QrRkFXZ3BRNzloMHJ4UUVBN1FhVGNDYXhiUG5RanNtYzd1eHN3M0YwWVdibE4wVEhFNGpTc1krNHZGNGtZaXJRbWNCbE5MSWJIUkhsb1RDQzJZcGh1ck5vNDZjSnFTR3I5Q1pieHZVbkt1U2tQelV5RnpkdEsyUTA2L0pFTXF1RzNMRC9TeGRvVld4RzgyVVhvYTRXaE54MXBRU1ZVdDREVU9rUzk1RVJsTkFVWnpteGczSjBhOC9PTGl2MHlXUm1EQmF5ZnU1WERJbUJaRGJ5WHVUdUQ1a2Y3UmhETjJXTkcwUVFWY0RPK1UyaEtTSWd6cGthL3lNTDVDUW9rL3pXQnM4MWJteUVnOG5vbXRwLy9KSkJocU5DOHpJcEIwWGhoMVMySHdaNCtDK3k1b2Ywd3Ardnl5SGp4RmZuWVR3K3hVRk83VDJjZTdwQWdNQkFBR2pnWTh3Z1l3d0h3WURWUjBqQkJnd0ZvQVVFb28zN1BselY5ZmtnN3ExYndQczNhNFJUczh3Q3dZRFZSMFBCQVFEQWdQNE1Gd0dDQ3NHQVFVRkJ3RUJCRkF3VGpBb0JnZ3JCZ0VGQlFjd0FvWWNhSFIwY0Rvdkx6RXlOeTR3TGpBdU1TOWpZV2x6YzNWbExtaDBiVEFpQmdnckJnRUZCUWN3QVlZV2FIUjBjRG92THpFeU55NHdMakF1TVRveU1EUTBNekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBWTZaaFdnb1IxbWxVT0JMbGJCRjBGZ0l2eFdxcFJld1lYUk1UT1NXTlF2K24xa1FkSFFOTGkrV2hhM21QRENMeHV5N2xCcEI3QTN5cDdvZDIwS3F6U25UcXFDNHExdkZqVGxpcmlxYVZvcjg4bXhpN2wrV2JIeW5VMTduYVpITG5DcEo1Q3RHdmJLREMvRGJGVHJvbDZ1SnBYTmp4eDNzZDBleUw5azNXYUIxWmgvZEV3cUtFdEw3K2xlSDg5MEM3enRTUGsyUkVjbFlLMTl1UWVPdzFNRzhpcVhBZTIrSlNzMUxuUzI3TFlsOG9xaWhDdkM2cGpTK3RkTThLTGk4cUd3MHBNM0lDc3pjQmlhOW5Od0tna0RIWkp1K1JyamtvZ3BGcjlWV3FJWG41cFMvekloc3Y2MEZaSnRKTXpsNnlPbXVYdVM4U1BtYk9Jd3FoblltQ2hnPT0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			"drk_sign":	"liUNzybvHigW6mbtnBEEf5Z9DnDywbYwFTkrmHCHLsVdDDvBJb4qoOmTVpmx3bJKRqngpLc0T0FKkUo1335Zkco76Mn22pPQND0rbg2RSEn3p8pigeQQl0o1QZklDRtVGs8iQENaYqpGAWPDOqPVXc1UIkGIOw7mVu6iNMuju9l-tRtvpMj055l6q42utR_FvjmsMJXO8shGaUyuvMOUU54ell8b05VndH08otaG0Ah-3PeN6a9udlXU0NNW3y8gqtjGgCeTitPWKWZuLoXXIfR6cCq0v6KF4USIEcrnme7sSvTicIWmKO15RdcouOgCfqVXPCDIrj--bFziY2rSVBBe9SxpN_4ZWKDqaBcix5wYr8IIQA53H5J1paoxgzh_y_M7ndim-r_rXGw0z4pHUMol2DDJGL5mAK2fcl8mvy-U8aq-KZCOS9Blk4uuI9ui2mFsrokDjangHmPxAfQZJrWdOi-qKSpvkRDuZNM1Llq5QBfu4BReNf2o_A8C16ImO1WhvuGm3N-nr6eqbxhmJkn_BO5tAN_ABqTDJWtzybdhT_fYbTad_YYcUbqZCclcRrn0z69PM10uR1P8E0dU9OTGQPT7JqCytawBU_4U2cYTelzKcFt0dVDGmaoKLhy7lZ0SNJNcQxNIRe-iSzl_aSaYs6lOMHD3OlsIEPvsJbk"
		},
		"payload":	{
			"version":	"TEE.RA.1.0",
			"timestamp":	1921869481084050,
			"scenario":	"sce_as_no_daa",
			"sign_alg":	"PS256",
			"hash_alg":	"HS256",
			"qta_img":	"mlBt4sRHz6uvqIQ61wqJtulDkxdt40iPL31xMKf9XEw",
			"qta_mem":	"m5yz5UDGYE1h-pwCB839HIO1I2WiZZoKuaxb87CcNl4",
			"tcb":	"",
			"ak_pub":	{
				"kty":	"RSA",
				"n":	"t5hBiQXxRErT3oEbnMWHzQmFRCPe5oa3wjmPLNOGWwNPLr2PUToSYKKU3KTaVDYGnTNeyZO5fTYnNgXi7zUmIxQ3avIsJn-esK5KGupd1fRp86nU472owIicKcbXibNTra-KIY090YDkhAVmqFQ5NFK_FtUuTXsRwTUcUWI_FXYaujnO1ml5kVS0UuWdPvD2vXdXRS39gWC5Vq3gLKFsoplfFB2-47nPMmSHdt2hsDpRFc6ShuRio34uOcPQHQ2Drh4ZUlCWtj4_QekEDjlSH31sGMM4xgZMyGJy5ilUU-JWLQ3zoGBBMSfJfnt4Bfic_iZDQUSRu3ZNV6yFJQuKGTYvaKjPWH4DQ1a85JCLONb6vUXIddfmKNJEz0pgLA_lRYpIVkD9ltVirOAG08x-ZcROZ5vPlq6KZSJLSr0in7PMKpoxo1y8DOp0Q9QxFk_o6o_RJGOCk9e_rqEGMRZma9Rit8Tu0aI-76AdLotuRlKJsmqfmyiy-wXh9Lyzbunh4Le8fK0Ci2cYWkti2l3mXVMdgGCjREk2IfuH6HEKCezj2FqswBnm-Y77BUw6855eKFPwMs9GejJFFgvyKgXyCD3t1JKtlykf0RCr05oHlgALa5soqhUOm3a1GuuPNhmcQfQ4tdLZjpN-ClZ1zTgqm5u9CR27KWIrpRXcHRjwDQc",
				"e":	"AQAB"
			}
		},
		"handler":	"provisioning-output"
	}`)

	_, _, akpub, err := GetDataFromAKCertNoDAA(jsonData)

	if err != nil {
		t.Errorf("generate nodaa scenario ak cert error %v", err)
	}
	t.Logf("akpub: %s", akpub)
}

/*
func TestOther1(t *testing.T) {
	//str := "aHR0cHM6Ly93d3cuZ2Vla3Nmb3JnZWVrcy5vcmcv"
	//str := "mlBt4sRHz6uvqIQ61wqJtulDkxdt40iPL31xMKf9XEw"
	//str := "igARsxqL_eiDN8fsHeGJHo_HPOEY4fQqfgONJF4hHOA"
	str := "IJFBNwUbczrU3mPnE4HLtHujR5-Op2O-UVSZ3l-HwJE"
	//strstd := "IJFBNwUbczrU3mPnE4HLtHujR5+Op2O+UVSZ3l+HwJE="
	//str := "liUNzybvHigW6mbtnBEEf5Z9DnDywbYwFTkrmHCHLsVdDDvBJb4qoOmTVpmx3bJKRqngpLc0T0FKkUo1335Zkco76Mn22pPQND0rbg2RSEn3p8pigeQQl0o1QZklDRtVGs8iQENaYqpGAWPDOqPVXc1UIkGIOw7mVu6iNMuju9l-tRtvpMj055l6q42utR_FvjmsMJXO8shGaUyuvMOUU54ell8b05VndH08otaG0Ah-3PeN6a9udlXU0NNW3y8gqtjGgCeTitPWKWZuLoXXIfR6cCq0v6KF4USIEcrnme7sSvTicIWmKO15RdcouOgCfqVXPCDIrj--bFziY2rSVBBe9SxpN_4ZWKDqaBcix5wYr8IIQA53H5J1paoxgzh_y_M7ndim-r_rXGw0z4pHUMol2DDJGL5mAK2fcl8mvy-U8aq-KZCOS9Blk4uuI9ui2mFsrokDjangHmPxAfQZJrWdOi-qKSpvkRDuZNM1Llq5QBfu4BReNf2o_A8C16ImO1WhvuGm3N-nr6eqbxhmJkn_BO5tAN_ABqTDJWtzybdhT_fYbTad_YYcUbqZCclcRrn0z69PM10uR1P8E0dU9OTGQPT7JqCytawBU_4U2cYTelzKcFt0dVDGmaoKLhy7lZ0SNJNcQxNIRe-iSzl_aSaYs6lOMHD3OlsIEPvsJbk"
	//str := "TUlJRWtqQ0NBM3FnQXdJQkFnSVJFV0ltOHVOTE5NU3ptZU9yRkhlM0p6Y3dEUVlKS29aSWh2Y05BUUVMQlFBd1BURUxNQWtHQTFVRUJoTUNRMDR4RHpBTkJnTlZCQW9UQmtoMVlYZGxhVEVkTUJzR0ExVUVBeE1VU0hWaGQyVnBJRWxVSUZCeWIyUjFZM1FnUTBFd0hoY05Nakl3TnpBNE1Ea3dOVFE1V2hjTk16Y3dOekEwTURrd05UUTVXakE2TVFzd0NRWURWUVFHRXdKRFRqRVBNQTBHQTFVRUNoTUdTSFZoZDJWcE1Sb3dHQVlEVlFRREV4RXdNalpRVUZZeE1FdEJNREExTVRNNVZEQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQQURDQ0Fnb0NnZ0lCQUxrSC82NkNTYU52MXhFVEd0WUtHZkNNbzVSTHlXYUdSVmFUbjR2QkhTUDk0ZXhvbU5DMjZSaHV4bEtVK1dNOEw0TjFSVDkxUnZRbVlUNUw3cG1xTzFlb043ei82YUhzZG5wZ3R4RUNSTzZ5WkxGcUs2UEhXd0twUnFkczRlYVcyZlZHVXF2MWZSRmVTaEtGM29YUVpYbmpLN1owS3M0QTZ6TmFNUnBmV3VNdlYxdHBXOFZVQVk0WG5lQjl0Zzk4d3ZRZUorR3k2cVoyVEFYVEkrekhUWnVMNENtTnUrTDdJTkNaT0VWUnNwVm1HVFhtSkhzbjhYQ0dWYllRMk5RT0FkdW9WUnRWaFZSMjk4SHQ0NUc1ejNTQzZtdEY1dTdPeXYvVlFQbytuS3hmRUNxLytjTUpIK3lsQWFCcXlya1hiYkU2SDRieU1PZUN5REFLSzlvSDdsdXNBenA2aURpN1QrRkFXZ3BRNzloMHJ4UUVBN1FhVGNDYXhiUG5RanNtYzd1eHN3M0YwWVdibE4wVEhFNGpTc1krNHZGNGtZaXJRbWNCbE5MSWJIUkhsb1RDQzJZcGh1ck5vNDZjSnFTR3I5Q1pieHZVbkt1U2tQelV5RnpkdEsyUTA2L0pFTXF1RzNMRC9TeGRvVld4RzgyVVhvYTRXaE54MXBRU1ZVdDREVU9rUzk1RVJsTkFVWnpteGczSjBhOC9PTGl2MHlXUm1EQmF5ZnU1WERJbUJaRGJ5WHVUdUQ1a2Y3UmhETjJXTkcwUVFWY0RPK1UyaEtTSWd6cGthL3lNTDVDUW9rL3pXQnM4MWJteUVnOG5vbXRwLy9KSkJocU5DOHpJcEIwWGhoMVMySHdaNCtDK3k1b2Ywd3Ardnl5SGp4RmZuWVR3K3hVRk83VDJjZTdwQWdNQkFBR2pnWTh3Z1l3d0h3WURWUjBqQkJnd0ZvQVVFb28zN1BselY5ZmtnN3ExYndQczNhNFJUczh3Q3dZRFZSMFBCQVFEQWdQNE1Gd0dDQ3NHQVFVRkJ3RUJCRkF3VGpBb0JnZ3JCZ0VGQlFjd0FvWWNhSFIwY0Rvdkx6RXlOeTR3TGpBdU1TOWpZV2x6YzNWbExtaDBiVEFpQmdnckJnRUZCUWN3QVlZV2FIUjBjRG92THpFeU55NHdMakF1TVRveU1EUTBNekFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBWTZaaFdnb1IxbWxVT0JMbGJCRjBGZ0l2eFdxcFJld1lYUk1UT1NXTlF2K24xa1FkSFFOTGkrV2hhM21QRENMeHV5N2xCcEI3QTN5cDdvZDIwS3F6U25UcXFDNHExdkZqVGxpcmlxYVZvcjg4bXhpN2wrV2JIeW5VMTduYVpITG5DcEo1Q3RHdmJLREMvRGJGVHJvbDZ1SnBYTmp4eDNzZDBleUw5azNXYUIxWmgvZEV3cUtFdEw3K2xlSDg5MEM3enRTUGsyUkVjbFlLMTl1UWVPdzFNRzhpcVhBZTIrSlNzMUxuUzI3TFlsOG9xaWhDdkM2cGpTK3RkTThLTGk4cUd3MHBNM0lDc3pjQmlhOW5Od0tna0RIWkp1K1JyamtvZ3BGcjlWV3FJWG41cFMvekloc3Y2MEZaSnRKTXpsNnlPbXVYdVM4U1BtYk9Jd3FoblltQ2hnPT0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	plain := make([]byte, 43)
	bytes := []byte(str)
	_, err := base64.RawURLEncoding.Decode(plain, bytes)
	//_, err := base64.StdEncoding.Decode(plain, bytes)
	//bytes, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		t.Errorf("Decode device cert error, %v", err)
	}
	t.Logf("len:%d, plaintext: %x", len(str), plain)
}
*/
