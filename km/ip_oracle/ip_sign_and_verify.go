package ip_oracle

import (
  "os"
  "time"
  "strings"
  "encoding/json"

  "github.com/aws/aws-sdk-go/service/kms"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  jwt "github.com/dgrijalva/jwt-go"
)

const AWS_DEFAULT_REGION = "ap-southeast-2"

type SigningMethodKMS struct {
  KeyId string
  GrantToken string
  SigningAlgorithm string
}

type IPOracleClaims struct {
    SourceIp string `json:"source_ip,omitempty"`
    jwt.StandardClaims
}

func NewSigningMethodKMS(keyId string) *SigningMethodKMS {
  return &SigningMethodKMS{
      KeyId: keyId,
      GrantToken: "CreateGrant",
      SigningAlgorithm: "RSASSA_PKCS1_V1_5_SHA_256",
  }
}

func (sm *SigningMethodKMS) Sign(signingString string, key interface{}) (string, error) {
  region := os.Getenv("AWS_REGION")
  if region == "" {
    region = AWS_DEFAULT_REGION
  }
  sess := session.Must(session.NewSession(&aws.Config{
    Region: aws.String(region),
  }))
  svc := kms.New(sess)

  keyId := sm.KeyId
  grantToken := sm.GrantToken
  signingAlgorithm := sm.SigningAlgorithm
  signInput := kms.SignInput{
    GrantTokens: []*string{&grantToken},
    KeyId: &keyId,
    Message: []byte(signingString),
    SigningAlgorithm: &signingAlgorithm,
  }

  signOutput, err := svc.Sign(&signInput)
  if err != nil {
    return "", err
  }

  sig := signOutput.Signature
  return jwt.EncodeSegment(sig), nil
}

func (sm *SigningMethodKMS) Alg() string {
  return "RS256"
}

func (sm *SigningMethodKMS) Verify(signingString string, signature string, key interface{}) error {
  region := os.Getenv("AWS_REGION")
  if region == "" {
    region = AWS_DEFAULT_REGION
  }
  sess := session.Must(session.NewSession(&aws.Config{
    Region: aws.String(region),
  }))
  svc := kms.New(sess)

  grantToken := sm.GrantToken
  keyId := sm.KeyId
  signingAlgorithm := sm.SigningAlgorithm
  sig, err := jwt.DecodeSegment(signature)
  verifyInput := kms.VerifyInput{
    GrantTokens: []*string{&grantToken},
    KeyId: &keyId,
    Message: []byte(signingString),
    Signature: sig,
    SigningAlgorithm: &signingAlgorithm,
  }

  _, err = svc.Verify(&verifyInput)
  return err
}

func MakeIPJWT(sourceIp string, sm jwt.SigningMethod) (string, error) {
  issuedAt := time.Now().UTC()
  expiresAt := issuedAt.Add(time.Minute * 10)
  
  claims := IPOracleClaims{
    sourceIp,
    jwt.StandardClaims{
      IssuedAt: issuedAt.UnixNano() / int64(time.Second),
      ExpiresAt: expiresAt.UnixNano() / int64(time.Second),
    },
  }
  token := jwt.NewWithClaims(sm, claims)
  signedString, err := token.SignedString(sm)
  return signedString, err
}

func ParseSignature(signedString string) (string, string) {
  signedString_parts := strings.Split(signedString, ".")
  return strings.Join(signedString_parts[0:2], "."), signedString_parts[2]
}

func VerifyIPJWT(signedString string, sm jwt.SigningMethod) (string, error) {
  signingString, signature := ParseSignature(signedString)
  err := sm.Verify(signingString, signature, nil)
  if err != nil {
    return "", err
  }
  return VerifyIPOracleClaims(signingString)
}

func VerifyIPOracleClaims(signingString string) (string, error) {
  signingString_parts := strings.Split(signingString, ".")
  claims, err := jwt.DecodeSegment(signingString_parts[1])
  if err != nil {
    return "", err
  }
  var ipOracleClaims IPOracleClaims
  json.Unmarshal(claims, &ipOracleClaims)
  return ipOracleClaims.SourceIp, ipOracleClaims.Valid()
}
