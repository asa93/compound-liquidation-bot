module gitlab.com/q-dev/exchange-rate-oracle

go 1.15

require (
	github.com/ethereum/go-ethereum v1.12.0
	github.com/go-kit/kit v0.9.0
	github.com/karalabe/usb v0.0.2 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	gitlab.com/q-dev/q-client v1.9.22-0.20211124080536-fe063185527d
	gitlab.com/q-dev/system-contracts v1.0.0-rc.5.0.20221004214545-578f7bdd1330
)

replace gitlab.com/q-dev/q-client v1.9.22-0.20210902222014-3ed08c979b9f => gitlab.com/q-dev/q-client v1.1.3-0.20221003065502-32b4e6c485df
