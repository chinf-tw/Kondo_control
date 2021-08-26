package main

import (
	"flag"
	"kondocontrol/internal/eeprom"
	"kondocontrol/internal/khr_3hv"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
)

var robot khr_3hv.RobotNum

func main() {
	var (
		lp = flag.String("left-port", "", "left port")
		rp = flag.String("right-port", "", "right port")
	)
	flag.Parse()
	if *lp == "" || *rp == "" {
		log.Fatalf("left and right port should not be empty, (lp: %s,rp: %s)", *lp, *rp)
	}
	// Set up leftOptions.
	leftOptions := serial.OpenOptions{
		PortName:          *lp,
		BaudRate:          1250000,
		DataBits:          8,
		StopBits:          1,
		MinimumReadSize:   3,
		ParityMode:        serial.PARITY_EVEN,
		RTSCTSFlowControl: false,
	}
	rightOptions := leftOptions
	rightOptions.PortName = *rp
	// Open the port.
	rightPort, err := serial.Open(rightOptions)
	if err != nil {
		log.Fatalf("rightPort.Open: %v", err)
	}
	leftPort, err := serial.Open(leftOptions)
	if err != nil {
		log.Fatalf("leftPort.Open: %v", err)
	}

	// Make sure to close it later.
	defer rightPort.Close()
	defer leftPort.Close()

	// init robot
	robot, err = khr_3hv.DefaultRobotNum(leftPort, rightPort)
	if err != nil {
		log.Fatal(err)
	}

	// run api
	apiRouter().Run(":8080")
}
func apiRouter() *gin.Engine {
	var api = gin.Default()

	api.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		// c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		// c.Header("Access-Control-Allow-Credentials", "true")

		// 別改
		c.Next()
	})
	api.GET("/wscontrol", wscontrol)
	api.GET("/batch_wscontrol", batchWscontrol)
	api.GET("/control", control)

	return api
}
func batchWscontrol(c *gin.Context) {
	var upgrader = websocket.Upgrader{} // use default options
	// websocket Upgrade
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wesocket upgrade failed"})
		return
	}
	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if msgType == websocket.TextMessage {
			// <number>,<angle>,<number>,<angle>,...
			cmd := strings.Split(string(msg), ",")
			if len(cmd)%2 != 0 || len(cmd) == 0 {
				log.Printf("不正確輸入: len(cmd) = %d\n%s\n", len(cmd), cmd)
				continue
			}
			for i := 0; i < len(cmd); i += 2 {
				if err := stringToPosition(cmd[i], cmd[i+1]); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}
func wscontrol(c *gin.Context) {
	var upgrader = websocket.Upgrader{} // use default options
	// websocket Upgrade
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wesocket upgrade failed"})
		return
	}
	defer ws.Close()
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if msgType == websocket.TextMessage {
			// <number>,<angle>
			cmd := strings.SplitN(string(msg), ",", 2)
			if err := stringToPosition(cmd[0], cmd[1]); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func control(c *gin.Context) {
	number := c.Query("number")
	angle := c.Query("angle")
	log.Println(number, angle)
	if err := stringToPosition(number, angle); err != nil {
		log.Println(err)
	}
}

var lastAngle [22]uint

func stringToPosition(number, angle string) error {

	num, err := strconv.Atoi(number)
	if err != nil {
		return errors.Wrap(err, "number")
	}
	ang, err := strconv.ParseUint(angle, 10, 16)
	if err != nil {
		return errors.Wrap(err, "angle")
	}
	if num > khr_3hv.LimitNum() {
		return errors.New("cmd[0] is bigger than LimitNum")
	}
	if uint16(ang) > eeprom.MaximumPosition || uint16(ang) < eeprom.MinimumPosition {
		return errors.New("angle is bigger than eeprom.MaximumPosition or smaller than eeprom.MinimumPosition")
	}
	if math.Abs(float64(lastAngle[num]-uint(ang))) < 50 {
		return nil
	}
	robot[num].SetPosition(uint(ang))
	lastAngle[num] = uint(ang)
	log.Printf("number: %s, angle: %s", number, angle)
	return nil
}
