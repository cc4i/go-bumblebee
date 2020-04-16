export interface Air {
    indexCityVHash: string;
	indexCity?:      string;
	stationIndex?:   number;
	aqi?:            number;
	city:           string;
	cityCN:        string;
	latitude?:       string;
	longitude?:      string;
	co?:             string;
	h?:              string;
	no2?:             string;
	o3?:              string;
	p?:               string;
	pm10?:            string;
	pm25?:            string;
	so2?:             string;
	t?:              string;
	w?:               string;
	s?:               string;  //Local measurement time
	tz?:              string; //Station timezone
	v?:               number;
}