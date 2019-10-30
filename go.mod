module siot-mqtt

go 1.12

replace (
	golang.org/x/crypto v0.0.0-20181127143415-eb0de9b17e85 => github.com/golang/crypto v0.0.0-20181127143415-eb0de9b17e85
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => github.com/golang/net v0.0.0-20181114220301-adae6a3d119a
)

require (
	github.com/astaxie/beego v1.11.1
	github.com/robfig/cron v1.1.0
	github.com/souliot/femq v0.0.0-20190512084014-a97c19731a9b
	github.com/souliot/fetcp v0.0.0-20180417075230-35effe73bb29
	github.com/souliot/siot-util v0.0.0-20190505091544-78536ad1a088
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94 // indirect
)
