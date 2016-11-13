package main

import (
  
   "math"
   
    
  
)



func Forecast( hpa float64, month int, wind float64, trend int, baro_top float64, baro_bottom float64) {
 forecast := [...]string{  "Settled fine", "Fine weather", "Becoming fine", "Fine, becoming less settled", "Fine, possible showers", "Fairly fine, improving", "Fairly fine, possible showers early", "Fairly fine, showery later", "Showery early, improving", "Changeable, mending", "Fairly fine, showers likely", "Rather unsettled clearing later", "Unsettled, probably improving", "Showery, bright intervals", "Showery, becoming less settled", "Changeable, some rain", "Unsettled, short fine intervals", "Unsettled, rain later", "Unsettled, some rain", "Mostly very unsettled", "Occasional rain, worsening", "Rain at times, very unsettled", "Rain at frequent intervals", "Rain, very unsettled", "Stormy, may improve", "Stormy, much rain" }
 
  rise_options  := [...]int{25,25,25,24,24,19,16,12,11,9,8,6,5,2,1,1,0,0,0,0,0,0}  
  steady_options := [...]int{25,25,25,25,25,25,23,23,22,18,15,13,10,4,1,1,0,0,0,0,0,0} 
  fall_options := []int {25,25,25,25,25,25,25,25,23,23,21,20,17,14,7,3,1,1,1,0,0,0} 

  var baro_range float64 = baro_top - baro_bottom
  var constant float64 = round((baro_range / 22))
  var season bool = ((month >= 4) && (month <= 9))
  
  if(wind == 0){
    hpa +=  (6 / 100) * baro_range;
  } else if (wind == 22 ) { // nne
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 45 ) { // ne
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) { // ene
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) { 
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) {
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) {
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) {
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) {
    hpa +=  (5 / 100) * baro_range;
  } else if (wind == 22 ) {
    hpa +=  (5 / 100) * baro_range;
  }
}

func round(f float64) float64 {
    return math.Floor(f + .5)
}
