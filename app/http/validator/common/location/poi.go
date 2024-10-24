package location

type Poi struct {
	Laitude   float64 `form:"latitude" json:"latitude"`
	Longitude float64 `form:"longitude" json:"longitude"`
}
