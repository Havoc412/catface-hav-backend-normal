package register_validator

import (
	"catface/app/core/container"
	"catface/app/global/consts"
	"catface/app/http/validator/common/upload_files"
	"catface/app/http/validator/common/websocket"
	"catface/app/http/validator/web/animal"
	"catface/app/http/validator/web/animal_like"
	"catface/app/http/validator/web/encounter"
	"catface/app/http/validator/web/encounter_like"
	"catface/app/http/validator/web/knowledge"
	"catface/app/http/validator/web/nlp"
	"catface/app/http/validator/web/search"
	"catface/app/http/validator/web/users"
)

// 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func WebRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string
	// Users 模块表单验证器按照 key => value 形式注册在容器，方便路由模块中调用
	key = consts.ValidatorPrefix + "UsersRegister"
	containers.Set(key, users.Register{})
	key = consts.ValidatorPrefix + "UsersLogin"
	containers.Set(key, users.Login{})

	key = consts.ValidatorPrefix + "UsersWeixinLogin"
	containers.Set(key, users.WeixinLogin{})

	key = consts.ValidatorPrefix + "RefreshToken"
	containers.Set(key, users.RefreshToken{})

	// Users基本操作（CURD）
	key = consts.ValidatorPrefix + "UsersShow"
	containers.Set(key, users.Show{})
	key = consts.ValidatorPrefix + "UsersStore"
	containers.Set(key, users.Store{})
	key = consts.ValidatorPrefix + "UsersUpdate"
	containers.Set(key, users.Update{})
	key = consts.ValidatorPrefix + "UsersDestroy"
	containers.Set(key, users.Destroy{})

	// 文件上传
	key = consts.ValidatorPrefix + "UploadFiles"
	containers.Set(key, upload_files.UpFiles{})

	// Websocket 连接验证器
	key = consts.ValidatorPrefix + "WebsocketConnect"
	containers.Set(key, websocket.Connect{})

	// Tag Animal
	key = consts.ValidatorPrefix + "AnimalList"
	containers.Set(key, animal.List{})
	key = consts.ValidatorPrefix + "AnimalDetail"
	containers.Set(key, animal.Detail{})
	key = consts.ValidatorPrefix + "AnimalCreate"
	containers.Set(key, animal.Create{})
	key = consts.ValidatorPrefix + "AnimalName"
	containers.Set(key, animal.Name{})

	key = consts.ValidatorPrefix + "AnimalLikeCreate"
	containers.Set(key, animal_like.Create{})
	key = consts.ValidatorPrefix + "AnimalLikeDelete"
	containers.Set(key, animal_like.Delete{})

	// TAG Encounter
	key = consts.ValidatorPrefix + "EncounterStore"
	containers.Set(key, encounter.Create{})
	key = consts.ValidatorPrefix + "EncounterList"
	containers.Set(key, encounter.List{})
	key = consts.ValidatorPrefix + "EncounterDetail"
	containers.Set(key, encounter.Detial{})

	key = consts.ValidatorPrefix + "EncounterLikeCreate"
	containers.Set(key, encounter_like.Create{})
	key = consts.ValidatorPrefix + "EncounterLikeDelete"
	containers.Set(key, encounter_like.Delete{})

	// TAG NLP
	key = consts.ValidatorPrefix + "NlpTitle"
	containers.Set(key, nlp.Title{})

	// TAG Search
	key = consts.ValidatorPrefix + "SearchAll"
	containers.Set(key, search.SearchAll{})

	// TAG Knowledge
	key = consts.ValidatorPrefix + "KnowledgeRandomList"
	containers.Set(key, knowledge.Random{})
}
