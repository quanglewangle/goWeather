package svgo1

import (
   "db2"
//   "os"
//   "xmldb"
   "math"
   "strconv"
   "bytes"
   "fmt"
   "github.com/ajstarks/svgo"
   "time"
   "wx"
  
)
const DegToRad = 0.017453292519943295769236907684886127134428718885417 
const RadToDeg = 57.295779513082320876798154814105170332405472466564   
	
const textStyle string =  "font-style:normal; font-variant:normal;font-weight:normal;font-stretch:normal;font-family:Arial;text-align:center;text-anchor:middle;stroke:none;fill-opacity:1; ";
const textReadingStyle string = "font-size:18px;font-weight: bold;fill: #ff3333;	fill-opacity:0.9"
const textUnitStyle string = "font-size:16px;font-weight: normal;fill: #ff3333; fill-opacity:0.9"
const lableStyle string = 	"font-size:18px;fill:#dcff85;font-weight:bold"
const histLableStyle string = 	"font-size:18px;fill:#222222;"

const circleStyle string = 	  "fill:none;stroke:black"
const historyStyle string = "fill:white;stroke:#303030;stroke-width:5;"
const rainDayHistoryLineStyle string = "fill: none;	fill-opacity:1;	stroke: #888800;stroke-width:1.5;stroke-opacity:1;"
const rainInstantHistoryLineStyle string = "fill: none;	fill-opacity:1;	stroke: #2222FF;stroke-width:1;stroke-opacity:1;"
const historyLineStyle string = "fill: none;	fill-opacity:1;	stroke: #2222FF;stroke-width:0.5;stroke-opacity:1;"
const pressureHistoryLineStyle string = "fill: none;	fill-opacity:1;	stroke: #2222FF;stroke-width:1.5;stroke-opacity:1;"


const cardinal_point_style string = "font-size:23px;font-weight: bold;	fill: #dcff85"
const intermediate_point_style string = "font-size:16px;font-weight: bold;	fill:  #dcff85"
const pointerShaftSyle string = "fill: #dcff85;	fill-opacity:1;	stroke: #303030;stroke-width:0.5;stroke-opacity:1;"
const gridStyle string = "stroke: #30FF33; stroke-width:0.5; fill: none"
const gridInsideWidth float64 = 595
const gridInsideHeight float64 = 295
const gridXOriginOffset float64 = 30
const historyUseableWidth = gridInsideWidth-gridXOriginOffset
const pointerCircleStyle string = "opacity:1; fill:#303030; fill-opacity:1; fill-rule:nonzero; stroke:black;"
const pointerShaftInnerStyle string = "fill: #303030;stroke:black;stroke-width:0.3;stroke-opacity:1;"
const bgRectStyle string = "fill:black; stroke: none;"
const bezelStyle string = "stroke-width:5.0; stroke: #505050"
const screwStyle string = "stroke-width:1.5; stroke: #AAAAAA"

const longDivStyle string = "fill:none;fill-rule:evenodd;stroke:#dcff85;stroke-width:2;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1;"
	
const shortDivStyle string = "fill:none;fill-rule:evenodd;stroke:#dcff85;stroke-width:0.5;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1;"
	
const size  		float64    = 300            	// The size of the image - must be divisible by 2
const cx    		float64    = size/2       	 	// x coordinate of the center
const cy   			float64    = cx             	// y coordinate of the center
const bezelR		int	   =    130
const screwOffset   int    =    22
const gridBoarderWidth float64 = 5

const pointerLength float64  	= size*0.38	            // Poninter length
const pointerWidth  float64   	= 6            	// width of the pointer
const inPointerLength float64 	= 30             	// length of the inner part of the shaft
const rPointerCircle  float64 	 = 20             	// Radius of the pointer inner circle

const rCardinalPoint float64  	= size*0.9             	// Radius of the circumference for the cardinal points
const xShadowOffset float64   	= 7	             	// x offset for the pointer shadow
const yShadowOffset float64   	= 7	             	// y offset for the pointer shadow

// Labels

const yVariable   float64     = 1.29*cy        // Position of the variable label (Wind)
const yCopyright  float64     = 0.83*cy        // Position of the text for the author
const yReading    float64     = 1.45*cy        // Position of the text for the speed
const yUnit       float64     = 1.7 *cy        // Position of the text for the units

const  nrHours int = 26

var cachedWindDir float64 = 0.0
var cachedWindSpeed float64  = 0.0
var cachedTemp float64 = 0.0
var cachedPress float64 = 990

var rExtDiv = int (math.Floor(size*0.42))
var angInc float64  

var rTxtScale = int (math.Floor(size*0.33))
var scaleVertOffset = 8 

func WindDirHistory() string {
   var aTime1, aTime2 time.Time
   var hist []wx.CurrentWeather
   var err error
  
   var x1, y1, x2, y2 int
   var r1, r2 float64 
   var sweepDir bool
   
   buf := new(bytes.Buffer)
   s := svg.New(buf)
   s.Start(int(size), int(size))

//   s.Circle(int(cx), int(cy), int(bezelR), bezelStyle)
   gunSight24Hrs(s)

   hist, _ = db.GetHistory()
   

   for i := 0; i<len(hist)-1; i++ {
       aTime1, err = time.Parse( time.RFC3339, hist[i].TimeStamp)
	   aTime2, err = time.Parse( time.RFC3339, hist[i+1].TimeStamp)
	   
       if(err != nil) {
		 sweepDir = false;
	     fmt.Println(err)
       } else {
//          fmt.Println( aTime1.Format(time.RFC3339))
		 
		  r1 = float64(timeToR(aTime1))
		  r2 = float64(timeToR(aTime2))
		 
		  x1 = int(round(r1 * math.Cos((hist[i].Wind.Direction-90.0)*DegToRad) + cx))
		  y1 = int(round(r1 * math.Sin((hist[i].Wind.Direction-90.0)*DegToRad) + cy))

		  x2 = int(round(r2 * math.Cos((hist[i+1].Wind.Direction-90.0)*DegToRad) + cx))
		  y2 = int(round(r2 * math.Sin((hist[i+1].Wind.Direction-90.0)*DegToRad) + cy))
		  
		  if((hist[i].Wind.Direction-90.0) > (hist[i+1].Wind.Direction-90.0)){
			sweepDir = false
		  } else {
			sweepDir = true
		  }
s.Arc(x1, y1, int(round(r1)), int(round(r2)), 0, false, sweepDir, x2, y2, historyLineStyle)
		//  s.Line(x1, y1, x2, y2, historyLineStyle)
	   }
	   
   }
   s.End()
   return string(buf.String())
}

func chooseInterval(max float64, min float64) int {
   if(max-min > 50.0){
	  return 10
   } else if (max-min > 9.0) {
	  return 5
   } else if (max-min > 4.0) {
	  return 2
   } 
   
   return 1
   
   
}

func RainHistoryGraph() string {
   var aTime1, aTime2 time.Time
   var hist []wx.CurrentWeather
   var err error
  // var hourNow = time.Time.Hour(time.Now())
   var x1, y1, x2, y2 int
   var max float64
   
   buf := new(bytes.Buffer)
   s := svg.New(buf)
  
   s.Start(int(size*2), int(size))
   
   s.Rect(int (gridBoarderWidth),int (gridBoarderWidth), int(gridInsideWidth)  , int (size - (gridBoarderWidth*2)), historyStyle)
  
  
   hist, max = db.GetRainHistory()
 
   if max < 1.0 {
	  max = 1.0
   }
   max = max * 1.05
   
   grid24Hrs(s, int(round(max)), 0, chooseInterval(max, 0))  
   for i := 0; i<len(hist)-1; i++ {
       aTime1, err = time.Parse( time.RFC3339, hist[i].TimeStamp)
	   aTime2, err = time.Parse( time.RFC3339, hist[i+1].TimeStamp)
	   
       if(err != nil) {
	       fmt.Println(err)
       } else {
		  x1 = timeToX(aTime1)
		  x2 = timeToX(aTime2)
		  y1 = valToY(float64(hist[i].Rain.Day), max, 0.0)
		  y2 = valToY(float64(hist[i+1].Rain.Day), max, 0.0)
//		  fmt.Println(x1, y1, x2, y2, float64(hist[i].Rain.Day), float64(hist[i+1].Rain.Day))
		  s.Line(x1, y1, x2, y2, rainDayHistoryLineStyle)

		  y1 = valToY(float64(hist[i].Rain.Instant), max, 0.0)
		  y2 = valToY(float64(hist[i+1].Rain.Instant), max, 0.0)
//		  fmt.Println(x1, y1, x2, y2, float64(hist[i].Rain.Day), float64(hist[i+1].Rain.Day))
		  s.Line(x1, y1, x2, y2, rainInstantHistoryLineStyle)


	   }
	   
   }
   
   s.End()
   return string(buf.String())
}

func HistoryGraph() string {
   var aTime1, aTime2 time.Time
   var hist []wx.CurrentWeather
   var err error
  // var hourNow = time.Time.Hour(time.Now())
   var x1, y1, x2, y2 int
  var max float64
  
   buf := new(bytes.Buffer)
   s := svg.New(buf)
  
   s.Start(int(size*2), int(size))
   
   s.Rect(int (gridBoarderWidth),int (gridBoarderWidth), int(gridInsideWidth)  , int (size - (gridBoarderWidth*2)), historyStyle)
   
   hist, max = db.GetHistory()
   if max < 5.0 {
	  max = 5.0
   }
   grid24Hrs(s, int(round(max)), 0, chooseInterval(max, 0))  
   
   for i := 0; i<len(hist)-1; i++ {
       aTime1, err = time.Parse( time.RFC3339, hist[i].TimeStamp)
	   aTime2, err = time.Parse( time.RFC3339, hist[i+1].TimeStamp)
	   
       if(err != nil) {
	       fmt.Println(err)
       } else {
//          fmt.Println( aTime1.Format(time.RFC3339))

		  x1 = timeToX(aTime1)
		  x2 = timeToX(aTime2)
//		  fmt.Println(aTime1, x1,  aTime2, x2)
	
		  y1 = valToY(float64(hist[i].Wind.Speed), max, 0.0)
		  y2 = valToY(float64(hist[i+1].Wind.Speed), max, 0.0)
		  s.Line(x1, y1, x2, y2, historyLineStyle)
	   }
	   
   }
   
   s.End()
   return string(buf.String())
}

func PressureHistoryGraph() string {
   var aTime1, aTime2 time.Time
   var hist []wx.CurrentWeather
   var err error
  // var hourNow = time.Time.Hour(time.Now())
   var x1, y1, x2, y2 int
  var max, min float64
  
   buf := new(bytes.Buffer)
   s := svg.New(buf)
  
   s.Start(int(size*2), int(size))
   
   s.Rect(int (gridBoarderWidth),int (gridBoarderWidth), int(gridInsideWidth)  , int (size - (gridBoarderWidth*2)), historyStyle)
   
   hist, max, min = db.GetPressureHistory()
   fmt.Println(max)
   if max < 5.0 {
	  max = 5.0
   }
   max = max * 1.01
   grid24Hrs(s, int(round(max)), int(round(min)), chooseInterval(max, min)) 
   
   for i := 0; i<len(hist)-1; i++ {
       aTime1, err = time.Parse( time.RFC3339, hist[i].TimeStamp)
	   aTime2, err = time.Parse( time.RFC3339, hist[i+1].TimeStamp)
	   
       if(err != nil) {
	       fmt.Println(err)
       } else {
//          fmt.Println( aTime1.Format(time.RFC3339))

		  x1 = timeToX(aTime1)
		  x2 = timeToX(aTime2)
	//	  fmt.Println(aTime1, x1,  aTime2, x2)
	
		  y1 = valToY(float64(hist[i].Pressure), max, 950.0)
		  y2 = valToY(float64(hist[i+1].Pressure), max, 950.0)
          //fmt.Println(y1, y2)
		  s.Line(x1, y1, x2, y2, pressureHistoryLineStyle)
	   }
	   
   }
    
   _, trend := (db.GetPressureAndTrend())
   s.Text(150, 50, strconv.Itoa(trend), histLableStyle)
   
   s.End()
   return string(buf.String())
}


func grid24Hrs(s *svg.SVG, maxY int, minY int, yScaleInterval int) {
   fmt.Println(maxY, minY, yScaleInterval)
   for i:=0; i < 24; i++ {
	  s.Line(hourToX(i),0, hourToX(i), int(size - (gridBoarderWidth*2))-25, gridStyle)
	  s.Text(hourToX(i), int(size - (gridBoarderWidth*2))-2, strconv.Itoa(i), histLableStyle)
   }
   
   // y = horizontal lines
   for y:=minY; y <= maxY; y = y + yScaleInterval {
	  s.Line(hourToX(0), valToY(float64(y), float64(maxY), float64(minY)), int(size*2), valToY(float64(y), float64(maxY), float64(minY)), gridStyle)
	  s.Text(hourToX(0)-10 , valToY(float64(y), float64(maxY), float64(minY))+10, strconv.Itoa(y), histLableStyle )
   }
   
}

func gunSight24Hrs(s *svg.SVG) {
   for i:=0; i <24; i=i+6 {
	  s.Circle(int(cx), int(cy), hourToR(i), gridStyle)
	  s.Text(hourToR(i)+int(cy), int(cx),  strconv.Itoa(i), histLableStyle)
   }
}
func valToY(val float64, maxVal float64, minVal float64) int {
   var pixPerVal = (gridInsideHeight-20.0)/(maxVal-minVal)
  // fmt.Println("valToY", gridInsideHeight, val, pixPerVal,round(gridInsideHeight-((val-minVal) * pixPerVal) ))
   return int(round(gridInsideHeight-((val-minVal) * pixPerVal)))-20
}

func hourToR(hour int) int {
   var rWidthOfHour float64 = (size/2.0)/24.0
   return int(round(rWidthOfHour * float64(hour)))
}

func timeToR(t time.Time) int {
   var rWidthOfMinute = (size/2.0)/(24.0*60.0)
   return (hourToR(time.Time.Hour(t) ) + int (round(rWidthOfMinute * float64(time.Time.Minute(t))))) 
}

func hourToX(hour int) int {   
   var xWidthOfHours float64 = historyUseableWidth/24.0
   var x int = int(round(gridXOriginOffset +  xWidthOfHours * float64(hour))) 

   return x

}

func timeToX(t time.Time) int {
   var xWidthOfMinute = historyUseableWidth/(24.0 * 60.0)
 
   return  hourToX(time.Time.Hour(t) ) + int(round(xWidthOfMinute * float64(time.Time.Minute(t))))   
}

func WindDial() string {
	  
   buf := new(bytes.Buffer)

  s := svg.New(buf)


//  s.Start(300, 300)
  s.Startview(300, 300, 0, 0, 300, 300)
    // var windDir = db.GetWind()
     var windDir = db.GetWind()
    
      compassDial(s, "Wind", windDir)
     
      pointer(s, cx, cy, windDir, cachedWindDir)
      cachedWindDir = windDir

      screw(s, screwOffset, screwOffset)
      screw(s, int(size)-screwOffset, screwOffset)
      screw(s, screwOffset, int(size)-screwOffset)
      screw(s, int(size)-screwOffset, int(size)-screwOffset)
     
  s.End()
  return string(buf.String())
}

func PressureDial() string {
   fmt.Printf("in PressureDial \n")
   buf := new(bytes.Buffer)
   s := svg.New(buf)

   s.Startview(300, 300, 0, 0, 300, 300)
   var  press, _ = db.GetPressureAndTrend()
   fmt.Printf("PressureDial pressure is %f\n", press)
    //s.Translate(100, 100)
   valueDial(s, press, 0, "Press", "hPa", 960, 1070, 1, 10, 1, -150, 150)
   
 //     fmt.Fprintf(os.Stdout, "%f \n", val2ang(float64(windSpeed), -150, 150, -10, 60))
     
   pointer(s, cx, cy, val2ang(press, -150, 150, 960, 1070),val2ang(cachedPress, -150, 150, 960, 1070))
   cachedPress = press

   screw(s, screwOffset, screwOffset)
   screw(s, int(size)-screwOffset, screwOffset)
   screw(s, screwOffset, int(size)-screwOffset)
   screw(s, int(size)-screwOffset, int(size)-screwOffset)
   s.End()
   return string(buf.String())
}

func TempDial() string {
	  
   buf := new(bytes.Buffer)

   s := svg.New(buf)

   s.Startview(300, 300, 0, 0, 300, 300)
   var temp = db.GetTemp()
  
   valueDial(s, temp, 0, "Temp", "C", -10, 30, 1, 10, 1, -150, 150)
   
 //     fmt.Fprintf(os.Stdout, "%f \n", val2ang(float64(windSpeed), -150, 150, -10, 60))
     
   pointer(s, cx, cy, val2ang(temp, -150, 150, -10, 30),val2ang(cachedTemp, -150, 150, -10, 30))
   cachedTemp = temp

   screw(s, screwOffset, screwOffset)
   screw(s, int(size)-screwOffset, screwOffset)
   screw(s, screwOffset, int(size)-screwOffset)
   screw(s, int(size)-screwOffset, int(size)-screwOffset)
     
   s.End()
   return string(buf.String())
}


func WindSpeedDial() string {
	  
   buf := new(bytes.Buffer)

  s := svg.New(buf)

  s.Startview(300, 300, 0, 0, 300, 300)
       var windSpeed = db.GetWindSpeed()
     // var windSpeed = xmldb.GetWindSpeed()
    //s.Translate(100, 100)
      valueDial(s, float64(windSpeed), 0, "Wind", "kt", 0, 60, 1, 10, 1, -150, 150)
   
 //     fmt.Fprintf(os.Stdout, "%f \n", val2ang(float64(windSpeed), -150, 150, 0, 60))
     
      pointer(s, cx, cy, val2ang(windSpeed, -150, 150, 0, 60),val2ang(cachedWindSpeed, -150, 150, 0, 60))
      cachedWindSpeed = windSpeed

      screw(s, screwOffset, screwOffset)
      screw(s, int(size)-screwOffset, screwOffset)
      screw(s, screwOffset, int(size)-screwOffset)
      screw(s, int(size)-screwOffset, int(size)-screwOffset)
     
  s.End()
  return string(buf.String())
}


func screw(s *svg.SVG, x int, y int){
	s.Circle(x, y, 7, screwStyle);
}

/* compassDial -- draw frame and dial. Write directions */
func compassDial(s *svg.SVG, unit string,  reading float64){
    s.Roundrect(0,0, int(size), int(size), 25, 25, bgRectStyle)
    s.Circle(int(cx), int(cy), int(bezelR), bezelStyle)
    s.Gstyle(textStyle)
    s.Text(int(cx), int( math.Floor(yVariable)), unit, textUnitStyle)
    s.Text(int(cx), int( math.Floor(yReading)), fmt.Sprintf("%03.f T", reading), textReadingStyle)

    s.Text(int(cx), int(math.Floor(rCardinalPoint)), "S", cardinal_point_style)
    s.Translate(int(cx), int(cy)) // fake rotate about a point with trans, rotate, trans -ve

   
       s.RotateTranslate(-int(cx), -int(cy), 90.0) 
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "W", cardinal_point_style)
       s.Gend()
 
       s.RotateTranslate(-int(cx), -int(cy), 180.0)
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "N", cardinal_point_style)
       s.Gend() 
       s.RotateTranslate(-int(cx), -int(cy), 270.0)
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "E", cardinal_point_style)
       s.Gend() 
       
       s.RotateTranslate(-int(cx), -int(cy), 45.0)
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "SW", intermediate_point_style)
       s.Gend()
       
        s.RotateTranslate(-int(cx), -int(cy), 135.0)
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "NW", intermediate_point_style)
       s.Gend()
       
        s.RotateTranslate(-int(cx), -int(cy), 225.0)
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "NE", intermediate_point_style)
       s.Gend()
       
        s.RotateTranslate(-int(cx), -int(cy), 315.0)
       s.Text(int(cx), int(math.Floor(rCardinalPoint)), "SE", intermediate_point_style)
       s.Gend()
       
       
     s.Gend() // Translate
    s.Gend() // Gstyle
 }
 
func valueDial(s *svg.SVG, reading float64, hiVal int, variable string, unit string, scaleLo int, scaleHi int, scaleInc int, scaleLong int, scaleMid int, angStart float64, angEnd float64){
 // hiVal is red pointer
 
    s.Roundrect(0,0, int(size), int(size), 25, 25, bgRectStyle)
    s.Circle(int(cx), int(cy), int(bezelR), bezelStyle)
    s.Gstyle(textStyle)
    s.Text(int(cx), int( math.Floor(yVariable)), variable, textUnitStyle)
    s.Text(int(cx), int( math.Floor(yReading)), fmt.Sprintf("%.1f %s", reading, unit), textReadingStyle)
      
  
    var rShortDiv = int(math.Floor(size*0.4))
    var rMidDiv =   int(math.Floor(size*0.38))
    var rLongDiv =  int(math.Floor(size*0.36))
    
    var nShortDiv int // number of short divs
    var nLongDiv int
   
    
    nShortDiv = (scaleHi - scaleLo) - scaleInc
    angInc  = (angEnd - angStart)/ float64(nShortDiv)
    nLongDiv = nShortDiv/scaleLong 			// number of lond divisions is number of shortDivs/how many short divs in a long
    
//    fmt.Fprintf(os.Stdout, "%d nShortDiv %d nLongDiv %f\n", nShortDiv, nLongDiv, angInc)
    
    var n int
    var m int
    var a float64
    var a1 float64
    
    // make longDiv ticks
  	for n=0; n<=nLongDiv; n++ {
		a = angStart + float64(n*scaleLong)*angInc;
//		 fmt.Fprintf(os.Stdout, "calling divPathString %f a  %d r1  %d r1\n", a, rLongDiv, rExtDiv)
   
		s.Path( divPathString(a,rLongDiv, rExtDiv)+"\n", longDivStyle)
		
		a1 = a;
		
		for m=1; m<scaleLong; m++ {
			a = a1 + float64(m)*angInc;
			  s.Path( divPathString(a,rShortDiv, rExtDiv)+"\n", shortDivStyle)
			
		}
		if(scaleMid==1){
			a = a1 + (float64(scaleLong/2)*angInc);
			s.Path( divPathString(a,rMidDiv, rExtDiv)+"\n", shortDivStyle)
		}
	}
  	a = angStart + float64(n*scaleLong)*angInc;
	s.Path( divPathString(a,rLongDiv, rExtDiv)+"\n", longDivStyle)
  	
  	var i int
    var numberInc int
    var x int
    var y int
    
    numberInc=scaleInc*scaleLong
//    fmt.Fprintf(os.Stdout, "doing scale, %d %d %d\n", scaleLo, scaleHi, numberInc)
	for  i= scaleLo;  i<= scaleHi; i+=numberInc {
		a = DegToRad*(val2ang( float64(i), angStart, angEnd, scaleLo, scaleHi));
		x = int(cx     +  float64(rTxtScale) * math.Sin(a));
		y = int(cy +      float64(scaleVertOffset) - float64(rTxtScale) * math.Cos(a));		
//		 fmt.Fprintf(os.Stdout, "calling text, %d %d %d\n", x, y, strconv.Itoa(i))
		s.Text(x, y, strconv.Itoa(i), lableStyle)
	
	}
    s.Gend()
}
  	
func divPathString (a float64, r1 int, r2 int) string {	

   
  //   fmt.Fprintf(os.Stdout, "in divPathString %f a  %d r1  %d r1\n", a, r1, r2)
   
	var x1 = int(cx + float64( r1) * math.Sin(a*DegToRad));
	var y1 = int(cy - float64( r1) * math.Cos(a*DegToRad));
	var x2 = int(cx + float64(rExtDiv) * math.Sin(a*DegToRad));
	var y2 = int(cy - float64(rExtDiv) * math.Cos(a*DegToRad));
	
	return fmt.Sprintf("M %d,%d L %d,%d", x1, y1, x2, y2)
	
	
}

/* pointer that rotates about cx cy. from north 0 - 360 */
func pointer(s *svg.SVG , fcx float64, fcy float64, rotateTo float64, rotateFrom float64){
    
    var tip int			= 8  // extent of pointy bit
    var halfWidth int   = 5
    
    // done all the positioning cacls now
    var cx int = int(math.Trunc(fcx))
    var cy int = int(math.Trunc(fcy))
    
    // make it take the shortest route
    if(math.Abs(rotateFrom-rotateTo) > 180.0) {
    	if(rotateFrom > 180.0){
    	    rotateFrom = toAntiClock(rotateFrom)
    	} else {
    	    rotateTo = toAntiClock(rotateTo)
   	    }
    }
    
    var animateTransformString string = fmt.Sprintf("<animateTransform attributeName=\"transform\" type=\"rotate\" from=\"%f 0 0\" to=\"%f 0 0\" dur=\"2s\"/>\n", 
    	rotateFrom, rotateTo)
    
    var pointerShaft string =  fmt.Sprintf("M 0,0  l %d,%d %d,%d %d,%d %d,%d %d,%d z", 
    
  	-halfWidth,
  	0, 
  	
  	0,
  	-int(pointerLength),
  	
  	halfWidth,
  	-tip,
  	
  	halfWidth,
  	tip,
  	
  	0,
  	int(pointerLength))

  	var pointerShaftInner string =  fmt.Sprintf("M 0,0  l %d,%d %d,%d %d,%d %d,%d %d,%d z", 
    
  	-halfWidth,
  	0, 
  	
  	0,
  	-int(pointerLength-70),
  	
  	halfWidth,
  	0,
  	
  	halfWidth,
  	0,
  	
  	0,
  	int(pointerLength-20))
  	
  	
  	// translate to pointer origin
  	s.Translate( int(cx), int(cy))

       s.Translate(0, 0) // fake rotate about a point with trans, rotate, trans -ve
         s.RotateTranslate(-0, 0, rotateTo)
         s.Path(pointerShaft, pointerShaftSyle)
         s.Path(pointerShaftInner, pointerShaftInnerStyle)
           
         s.Writer.Write([]byte(animateTransformString))

        s.Gend() 
     s.Gend()
     s.Gend() // origin translate
     s.Circle(int(cx), int(cy), int(rPointerCircle), pointerCircleStyle);
     
}

func toAntiClock(b float64) float64 {
	return -(360.0 - b)
}

func val2ang(v float64, angStart float64, angEnd float64, scaleLo int, scaleHi int) float64 {
/*
This function computes the rotation angle for the pointer.
*/	
    var scaleHighest = float64(scaleHi)
    var scaleLowest = float64(scaleLo)
    var a float64
    var angOrigin float64
    
    if(angStart < 0){
    	angOrigin = 360+angStart
    } else {
    	angOrigin = angStart
    }
   
    if (scaleHighest > scaleLowest){	
//    	fmt.Fprintf(os.Stdout, " in if %f %f %f \n", math.Abs(angEnd-angStart), (scaleHighest-scaleLowest), (v-scaleLowest))
		a = math.Abs(angEnd-angStart)/(scaleHighest-scaleLowest)*(v-scaleLowest)+angOrigin;
	}else{
		a = 0;
	}
//	 fmt.Fprintf(os.Stdout, "%f %f %f %f %f %f\n", angEnd, angStart, scaleHighest, scaleLowest, v, a)
	
	return a;
}

func round(f float64) float64 {
    return math.Floor(f + .5)
}
 
