package test

import (
	"catface/app/model_es"
	_ "catface/bootstrap"
	"testing"
)

var Know_vacc_1 = []model_es.Knowledge{
	{
		Dirs:    []string{"疫苗", "喵三联"},
		Title:   "概念",
		Content: "喵三联疫苗是一种针对猫咪的疫苗组合，通常包含三种疫苗，以提供对猫瘟、白血、传腹三种主要疾病的防护。",
	},
	{
		Dirs:    []string{"疫苗", "喵三联"},
		Title:   "接种周期",
		Content: "1. 幼猫：通常在6-8周大开始接种第一剂，之后每3-4周接种一次，直到完成基础免疫。\n2. 成猫：完成基础免疫后，每年进行一次加强接种。",
	},
	{
		Dirs:    []string{"疫苗", "喵三联"},
		Title:   "副作用",
		Content: "像所有疫苗一样，喵三联疫苗可能会有轻微的副作用，如注射部位的红肿、疼痛或轻微的发烧。",
	},
}

var Know_food = []model_es.Knowledge{
	{
		Dirs:    []string{"食物"},
		Title:   "猫粮",
		Content: "选择专为猫设计的猫粮，考虑其年龄、体重、健康状况和活动水平，确保营养均衡，避免含有对猫有害的成分。",
	},
	{
		Dirs:    []string{"食物"},
		Title:   "猫条",
		Content: "猫条可作为零食或训练工具，选择高质量猫条，适量喂食，注意热量控制，避免肥胖。",
	},
	{
		Dirs:    []string{"食物"},
		Title:   "人类食物对猫的影响",
		Content: "人类食物可能不适合猫，某些食物对猫有毒，高盐、高糖或高脂肪食物可能导致健康问题，不应作为猫咪日常饮食。",
	},
	{
		Dirs:    []string{"食物"},
		Title:   "猫粮的营养成分",
		Content: "猫粮应包含高质量的蛋白质、脂肪、维生素和矿物质，以满足猫的营养需求。",
	},
}

var Know_WHU = []model_es.Knowledge{
	{
		Dirs:    []string{"WHU"},
		Title:   "狂犬疫苗接种点",
		Content: "地质医院是离武大最近的接种点。",
	},
}

var Know_viuse = []model_es.Knowledge{
	{
		Dirs:    []string{"狂犬病"},
		Title:   "概念",
		Content: "狂犬病毒是一种通过动物咬伤或唾液传播的病毒，主要影响中枢神经系统，对人类和动物都具有致命性。",
	},
	{
		Dirs:    []string{"狂犬病"},
		Title:   "人类症状",
		Content: "人类感染狂犬病毒后，初期可能表现为发热、头痛、恶心等症状，随后可能出现恐惧、兴奋、幻觉、瘫痪和昏迷等严重症状。",
	},
	{
		Dirs:    []string{"狂犬病"},
		Title:   "传播途径",
		Content: "狂犬病毒主要通过被感染动物的咬伤传播，也可能通过唾液直接接触眼睛、鼻子或口腔而传播。",
	},
	{
		Dirs:    []string{"狂犬病"},
		Title:   "与猫",
		Content: "猫可以感染狂犬病毒，但通常不表现出明显症状，可能成为无症状携带者，传播风险相对较低。",
	},
	{
		Dirs:    []string{"狂犬病"},
		Title:   "猫狗预防",
		Content: "定期给猫接种狂犬病疫苗，避免与野生动物接触，及时处理伤口，并在疑似暴露后采取隔离措施。",
	},
	{
		Dirs:    []string{"狂犬病"},
		Title:   "人类预防",
		Content: "避免接触可能携带狂犬病毒的动物，被动物咬伤后立即清洗伤口并接种疫苗，了解并遵守当地狂犬病预防和动物管理的法律法规。",
	},
	{
		Dirs:    []string{"狂犬病"},
		Title:   "治疗方法",
		Content: "狂犬病目前没有特效治疗，治疗主要针对症状，因此预防措施至关重要。",
	},
}

func TestKnowledgesES(t *testing.T) {
	// 定义需要插入的数据
	knowledges := append((append(append(Know_viuse, Know_food...), Know_WHU...)), Know_vacc_1...)

	// 遍历并插入每个知识条目
	for _, knowledge := range knowledges {
		err := knowledge.InsertDocument()
		if err != nil {
			t.Fatalf("Error inserting document: %v", err)
		}
	}
}

func TestRandomSearch(t *testing.T) {
	knowledges, err := model_es.CreateKnowledgeESFactory().RandomDocuments(5)
	if err != nil {
		t.Fatalf("Error retrieving random documents: %v", err)
	}

	t.Log("随机搜索结果：", len(knowledges))
	for _, knowledge := range knowledges {
		t.Log(knowledge)
	}
}
