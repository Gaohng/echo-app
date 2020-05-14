package gsafe

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"strconv"
	"time"
)

type JwsHelperOpt struct {
	Audience       string
	Issuer         string
	Timeout        int64  //秒
	PublicKeyPath  string //接受者需要提供公钥(不需要提供公钥)
	PrivateKeyPath string //签发者需要知道私钥
	HashIdsSalt    string //配置后使用hashIds　混淆UserId
}

type JwsHelper struct {
	opt        JwsHelperOpt
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJwsHelper(opt JwsHelperOpt) (*JwsHelper, error) {
	jwsHelper := &JwsHelper{opt: opt}
	if opt.PrivateKeyPath != "" {
		keyData, _ := ioutil.ReadFile(opt.PrivateKeyPath)
		var err error
		jwsHelper.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			return nil, errors.Wrap(err, "ParseRSAPrivateKeyFromPEM")
		}
	}

	if opt.PublicKeyPath != "" {
		keyData, _ := ioutil.ReadFile(opt.PublicKeyPath)
		var err error
		jwsHelper.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(keyData)
		if err != nil {
			return nil, errors.Wrap(err, "ParseRSAPrivateKeyFromPEM")
		}
	}
	return jwsHelper, nil
}

func (jws *JwsHelper) CreateToken(userId int64, payload string) (string, error) {
	var userIdHash string
	var err error
	if jws.opt.HashIdsSalt != "" {
		userIdHash, err = EncodeInt64(userId, jws.opt.HashIdsSalt)
		if err != nil {
			return "", errors.Wrap(err, "hashids.EncodeInt64")
		}
	} else {
		userIdHash = strconv.FormatInt(userId, 10)
	}

	claims := jwt.StandardClaims{
		Audience:  jws.opt.Audience,
		ExpiresAt: time.Now().Unix() + jws.opt.Timeout,
		Id:        userIdHash,
		IssuedAt:  time.Now().Unix(),
		Issuer:    jws.opt.Issuer,
		Subject:   payload,
	}
	//log.Printf("claims : %+v", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)
	res, err := token.SignedString(jws.privateKey)
	return res, err
}

func (jws *JwsHelper) ParseToken(tokenStr string) (id int64, payload string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, err error) {
		return jws.publicKey, nil
	})
	if err != nil {
		return 0, "", errors.Wrap(err, "jwt.Parse")
	}
	if token.Claims.Valid() != nil {
		return 0, "", errors.Wrap(err, "token.Claims.Valid()")
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.Errorf("claim is not StandardClaims: %+v", token.Claims)
	}

	if !claim.VerifyAudience(jws.opt.Audience, true) {
		return 0, "", errors.Errorf("claim.VerifyAudience: %+v", claim)
	}

	if !claim.VerifyIssuer(jws.opt.Issuer, true) {
		return 0, "", errors.Errorf("claim.VerifyIssuer: %+v", claim)
	}

	userIdStr, ok := claim["jti"].(string)
	var userId int64
	if !ok {
		return 0, "", errors.Errorf("claim is not StandardClaims: %+v", token.Claims)
	}
	if jws.opt.HashIdsSalt != "" {
		userId, err = DecodeInt64(userIdStr, jws.opt.HashIdsSalt)
		if err != nil {
			return 0, "", errors.Wrap(err, "hashids.EncodeInt64")
		}
	} else {
		userId, err = strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			return 0, "", errors.Wrap(err, "strconv.ParseInt")
		}
	}

	payload, ok = claim["sub"].(string)
	if !ok {
		return 0, "", errors.Errorf("claim is not StandardClaims sub is not string: %+v", token.Claims)
	}
	return userId, payload, nil
}
