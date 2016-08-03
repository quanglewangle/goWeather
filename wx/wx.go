package wx
type Temperature struct {
    Unit string `xml:"Unit"`
    Value float64 `xml:"value"`
    FeelsLike float64 `xml:"feels_like"`
}

type CurrentRain struct {
    Instant float64
    Day float64
}

type CurrentWind struct {
    Speed float64 `xml:"speed"`
    Direction float64 `xml:"direction"`
    Gusts float64 `xml:"gusts"`
    
}

type CurrentWeather struct {
       TimeStamp string `xml:"timestamp"`
       Wind CurrentWind `xml:"wind"`
       Temp Temperature `xml:"temperature"`
       Rain CurrentRain
       Pressure float64
}

type HistoricWeather struct {
       CurrentWeathers []CurrentWeather `xml:"currentWeather"`
}
