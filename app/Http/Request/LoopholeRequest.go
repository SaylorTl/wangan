package Request

type LoopholeCheckValidate struct {
	//Poc_id    int    `form:"poc_id" to:"poc_id"  binding:"required"`
	Threat_id string `form:"threat_id" to:"threat_id" binding:"required"`
	//User_id   int    `form:"user_id" to:"user_id"  binding:"required"`
}
