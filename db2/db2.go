package db

import ( "database/sql"
         _ "github.com/go-sql-driver/mysql"
//	"strconv"
         "fmt"
         "wx"
)

var db *sql.DB
var err error

// idempotent Open
func OpenDatabase() {
    
    db, err = sql.Open("mysql", "demo:password@/weather")
    if err != nil {
       fmt.Println(err)
    }
}

func GetHistory() ([]wx.CurrentWeather, float64) {
       const query string = "SELECT " +
                        " DATE_FORMAT(time_stamp, '%Y-%m-%dT%T+00:00') , " +
                        " avg(wind_speed_kt), " +
                        " avg(wind_dir_degree) " +
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() " +
                        " GROUP BY " + 
                        " round(UNIX_TIMESTAMP(time_stamp) / 180) "
                        
       const maxQuery string = "SELECT " +
                        " max(wind_speed_kt) " +
       
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() " 
        var got []wx.CurrentWeather
        var maxWind float64
        
   // db, err := sql.Open("mysql", "demo:password@/weather")
  //  db.parseTime=true
    rows, err := db.Query(query)
    if err != nil {
       fmt.Println("query %v", err)
    }
   
    for rows.Next() {
            var r wx.CurrentWeather
            err = rows.Scan(&r.TimeStamp, &r.Wind.Speed, &r.Wind.Direction)
            if err != nil {
                    fmt.Println("Scan: %v", err)
            }
            got = append(got, r)
    }
    err = db.QueryRow(maxQuery).Scan(&maxWind)
    if err != nil {
	    fmt.Println(err)
    }
  //  fmt.Println(got)
    return got, maxWind
}

func GetPressureAndTrend() (float64, int) {
  var pressure1 float64
  var pressure2 float64
  
  const trendQuery1 string = "SELECT " +
                        " pressure_hpa" +
       
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = DATE(CURRENT_DATE) AND HOUR(time_stamp) = HOUR(CURRENT_TIME) -1 " +
                        " ORDER BY time_stamp desc " +
                        "LIMIT 1"
  
  const trendQuery2 string = "SELECT " +
                        " pressure_hpa" +
       
                        " FROM reading " +
                        " ORDER BY time_stamp desc " +
                        "LIMIT 1"
                        
  row := db.QueryRow(trendQuery1)
  err := row.Scan(&pressure1)
  if err != nil {
	    fmt.Println(err)
  }
  
  row = db.QueryRow(trendQuery2)
  err = row.Scan(&pressure2)
  if err != nil {
	    fmt.Println(err)
  }
  fmt.Println("GetPressureandTrend pressures", pressure1, pressure2)
  if(pressure1 > (pressure2 + 0.1)) {
       return   pressure2, -1
  } else if (pressure2 > (pressure1 + 0.1)) {
       return   pressure2, 1
  } else {
       return  pressure2, 0
  }
}
  
func GetPressureHistory() ([]wx.CurrentWeather, float64, float64) {
       const query string = "SELECT " +
                        " DATE_FORMAT(time_stamp, '%Y-%m-%dT%T+00:00') , " +             
                        
                        " pressure_hpa " +
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() "
                        
       const maxQuery string = "SELECT " +
                        " max(pressure_hpa)" +
 
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() "
                        
       const minQuery string = "SELECT " +
                        " min(pressure_hpa)" +
       
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() " 
        var got []wx.CurrentWeather
        var maxPressure float64
        var minPressure float64
        
   // db, err := sql.Open("mysql", "demo:password@/weather")
  //  db.parseTime=true
    rows, err := db.Query(query)
    if err != nil {
       fmt.Println("query %v", err)
    }
   
    for rows.Next() {
            var r wx.CurrentWeather
            err = rows.Scan(&r.TimeStamp, &r.Pressure)
            if err != nil {
                    fmt.Println("Scan: %v", err)
            }
            got = append(got, r)
    }
   err = db.QueryRow(maxQuery).Scan(&maxPressure)
    if err != nil {
	    fmt.Println(err)
    }
    err = db.QueryRow(minQuery).Scan(&minPressure)
    if err != nil {
	    fmt.Println(err)
    }
  //  fmt.Println(got)
    return got, maxPressure, minPressure
}

func GetRainHistory() ([]wx.CurrentWeather, float64) {
       const query string = "SELECT " +
                        " DATE_FORMAT(time_stamp, '%Y-%m-%dT%T+00:00') , " +
                        " max(rain_mm), " +
                        " max(rain_day_mm) " +
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() " +
                        " GROUP BY " + 
                        " round(UNIX_TIMESTAMP(time_stamp) / 360) "
                        
       const maxQuery string = "SELECT " +
                        " max(rain_mm), " +
                        " max(rain_day_mm) " +
                        " FROM reading " +
                        " WHERE DATE(time_stamp) = CURDATE() " 
                 
        var got []wx.CurrentWeather
        var maxDay, maxInstant float64
        
   // db, err := sql.Open("mysql", "demo:password@/weather")
  //  db.parseTime=true
    rows, err := db.Query(query)
    if err != nil {
     fmt.Println("query %v", err)
    }
   
    for rows.Next() {
            var r wx.CurrentWeather
            err = rows.Scan(&r.TimeStamp, &r.Rain.Instant, &r.Rain.Day)
            if err != nil {
                    fmt.Println("Scan: %v", err)
            }
            got = append(got, r)
    }
 
    err = db.QueryRow(maxQuery).Scan(&maxInstant, &maxDay)
    if err != nil {
	    fmt.Println(err)
    }
    if maxInstant > maxDay {
       return got, maxInstant
    } else {
       return got, maxDay
    }
}

func GetWind() float64 {
    // db, err := sql.Open("mysql", "demo:password@/weather")
    if err != nil {
        panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
    }    
   // defer db.Close()

    stmtOut, err := db.Prepare("SELECT wind_dir_degree FROM reading ORDER BY time_stamp DESC LIMIT 1")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var windDir float64 // we "scan" the result in here

    
    err = stmtOut.QueryRow().Scan(&windDir) // WHERE number = 13
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    fmt.Printf("GetWind: windDir is %f\n", windDir)

  //  retS, err = strconv.ParseFloat(windDir, 64)  
  //  if err != nil {
  //      panic(err.Error())
  //  }
    return windDir
}

func GetRain() (float64, float64) {
    // db, err := sql.Open("mysql", "demo:password@/weather")
    if err != nil {
        panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
    }    
   // defer db.Close()

    stmtOut, err := db.Prepare("SELECT rain_mm, rain_day_mm FROM reading ORDER BY time_stamp DESC LIMIT 1")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var rain_mm, rain_day_mm float64 // we "scan" the result in here

    
    err = stmtOut.QueryRow().Scan(&rain_mm, &rain_day_mm) // WHERE number = 13
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
  
  
    return rain_mm, rain_day_mm
}



func GetWindSpeed() float64 {
    //db, err := sql.Open("mysql", "demo:password@/weather")
    if err != nil {
        panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
    }    
  //  defer db.Close()

    stmtOut, err := db.Prepare("SELECT wind_speed_kt FROM reading ORDER BY time_stamp DESC LIMIT 1")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var windSpeed float64 // we "scan" the result in here

    
    err = stmtOut.QueryRow().Scan(&windSpeed) // WHERE number = 13
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    fmt.Printf("GetWindSpeed: windDir is %f\n", windSpeed)

 
    return windSpeed
}
func GetTemp() float64 {
    //db, err := sql.Open("mysql", "demo:password@/weather")
    if err != nil {
        panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
    }    
   // defer db.Close()

    stmtOut, err := db.Prepare("SELECT temp_c FROM reading ORDER BY time_stamp DESC LIMIT 1")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var temp float64 // we "scan" the result in here
   
    err = stmtOut.QueryRow().Scan(&temp) // WHERE number = 13
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    fmt.Printf("temp is %f\n", temp)

 
    return temp
}

func GetPressure() float64 {
   
    fmt.Printf("in GetPressure:")
    stmtOut, err := db.Prepare("SELECT pressure_hpa FROM reading ORDER BY time_stamp DESC LIMIT 1")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    defer stmtOut.Close()

    var press float64 // we "scan" the result in here

    
    err = stmtOut.QueryRow().Scan(&press) // WHERE number = 13
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    fmt.Printf("GetPressure: pressure is %f\n", press)

 
    return press
}
