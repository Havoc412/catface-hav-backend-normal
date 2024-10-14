package model

/**
 * @description: 保留 Top 3, 辅助 catface - breed 子模型判断； 单独建表，因为只会被 CatFace 模块使用。
 * @return {*}
 */
type AnmFaceBreed struct { // TODO 迁移 python 的时候再考虑一下细节
	BriefModel
	Top1  uint8
	Prob1 float64
	Top2  uint8
	Prob2 float64
	Top3  uint8
	Prob3 float64

	AnimalId int64 // INFO 外键设定?
	Animal   Animal
}
