package curd

import "catface/app/model"

func CreateEncounterLikeCurdFactory() *EncounterLikeCurd {
	return &EncounterLikeCurd{model.CreateEncounterLikeFactory("")}
}

type EncounterLikeCurd struct {
	encounterLike *model.EncounterLike
}
