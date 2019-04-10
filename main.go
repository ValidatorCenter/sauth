package main

import (
	"fmt"
	"net/http"
	"strings"

	// для API
	"gopkg.in/macaron.v1"

	// для авторизации/регистрации
	"github.com/miguelmota/go-ethereum-hdwallet"

	// новая seed-фраза
	"github.com/tyler-smith/go-bip39"
)

type RetJSONSeed struct {
	Status   bool   `json:"status"`
	Mnemonic string `json:"mnemonic"`
	ErrMsg   string `json:"err_msg"`
}

type RetJSONPriv struct {
	RetJSONSeed
	Address string `json:"address"`
	Privkey string `json:"priv_key"`
}

// Авторизация по Seed-фразе
func AuthMnemonic(seedPhr string) (string, string, error) {
	wallet, err := hdwallet.NewFromMnemonic(seedPhr)
	if err != nil {
		//panic(err)
		return "", "", err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", err
	}

	//M+`в нижнем регистре(без видущего нуля)`
	strAdrs := account.Address.String()                                    // 0x512B699Ab21542B8875609593e845818f301903B
	addrss := fmt.Sprintf("M%s", strings.ToLower(strAdrs[1:len(strAdrs)])) // Mx512b699ab21542b8875609593e845818f301903b
	privKeyStr, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return "", "", err
	}
	return addrss, privKeyStr, nil
}

// Генерация новой Seed-фразы
func NewMnemonic() string {
	// Создаёт мнемонику для запоминания или удобный для пользователя seed
	// Мнемоника: это seed фраза
	entropy, _ := bip39.NewEntropy(128) //biteSize должен быть кратен 32 и находиться в пределах включенного диапазона {128, 256}
	Mnemonic, _ := bip39.NewMnemonic(entropy)
	return Mnemonic
}

// Возвращает JSON новую seed-фразу
func hndAPINewMnemonic(ctx *macaron.Context) {
	retDt := RetJSONSeed{}

	// Создаёт мнемонику для запоминания или удобный для пользователя seed
	// Мнемоника: это seed фраза
	retDt.Mnemonic = NewMnemonic()
	retDt.Status = true // исполнен без ошибок
	retDt.ErrMsg = ""   // нет ошибок

	// возврат JSON данных
	ctx.JSON(200, &retDt)
}

// Регистрация и авторизация
func hndAPIAuthUser(ctx *macaron.Context) {
	var seedPhr string
	retDt := RetJSONPriv{}

	ctx.Req.ParseForm()
	ctx.Resp.WriteHeader(http.StatusOK)
	seedPhr = ctx.Req.PostFormValue("sp")
	if seedPhr == "" {
		seedPhr = ctx.Req.FormValue("sp")
	}

	if seedPhr != "" {
		addrss, privKeyStr, err := AuthMnemonic(seedPhr)
		if err != nil {
			retDt.Status = false
			retDt.ErrMsg = err.Error()

			// редирект
			ctx.JSON(200, &retDt)
			return
		}

		retDt.Privkey = privKeyStr
		retDt.Address = addrss
		retDt.Mnemonic = seedPhr
		retDt.Status = true // исполнен без ошибок
		retDt.ErrMsg = ""   // нет ошибок
	} else {
		retDt.Status = false
		retDt.ErrMsg = "No seed-phrase"
	}

	// возврат JSON данных
	ctx.JSON(200, &retDt)
}

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())

	m.Route("/api/v1/authSeed", "GET,POST", hndAPIAuthUser) // регистрация и авторизация
	m.Get("/api/v1/newMnemonic", hndAPINewMnemonic)         // новая seed-фраза, регистрация нового аккаунта в сети Minter

	m.Run(3999) //port:3999
}
