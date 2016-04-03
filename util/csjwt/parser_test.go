package csjwt_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"runtime"
	"runtime/debug"

	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/stretchr/testify/assert"
)

var (
	defaultKeyFunc csjwt.Keyfunc = func(t csjwt.Token) (csjwt.Key, error) {
		return csjwt.WithRSAPublicKeyFromFile("test/sample_key.pub"), nil
	}
	emptyKeyFunc csjwt.Keyfunc = func(t csjwt.Token) (csjwt.Key, error) { return csjwt.Key{}, nil }
	errorKeyFunc csjwt.Keyfunc = func(t csjwt.Token) (csjwt.Key, error) { return csjwt.Key{}, fmt.Errorf("error loading key") }
	nilKeyFunc   csjwt.Keyfunc = nil
)

var jwtTestData = []struct {
	name        string
	tokenString []byte
	keyfunc     csjwt.Keyfunc
	claims      csjwt.MapClaims
	valid       bool
	errors      uint32
	parser      *csjwt.Parser
}{
	{
		"basic",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		defaultKeyFunc,
		csjwt.MapClaims{"foo": "bar"},
		true,
		0,
		nil,
	},
	{
		"basic expired",
		nil, // autogen
		defaultKeyFunc,
		csjwt.MapClaims{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		csjwt.ValidationErrorExpired,
		nil,
	},
	{
		"basic nbf",
		nil, // autogen
		defaultKeyFunc,
		csjwt.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		csjwt.ValidationErrorNotValidYet,
		nil,
	},
	{
		"expired and nbf",
		nil, // autogen
		defaultKeyFunc,
		csjwt.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		csjwt.ValidationErrorNotValidYet | csjwt.ValidationErrorExpired,
		nil,
	},
	{
		"basic invalid",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		defaultKeyFunc,
		csjwt.MapClaims{"foo": "bar"},
		false,
		csjwt.ValidationErrorSignatureInvalid,
		nil,
	},
	{
		"basic nokeyfunc",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		nilKeyFunc,
		csjwt.MapClaims{"foo": "bar"},
		false,
		csjwt.ValidationErrorUnverifiable,
		nil,
	},
	{
		"basic nokey",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		emptyKeyFunc,
		csjwt.MapClaims{"foo": "bar"},
		false,
		csjwt.ValidationErrorSignatureInvalid,
		nil,
	},
	{
		"basic errorkey",
		[]byte("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg"),
		errorKeyFunc,
		csjwt.MapClaims{"foo": "bar"},
		false,
		csjwt.ValidationErrorUnverifiable,
		nil,
	},
	{
		"invalid signing method",
		nil,
		defaultKeyFunc,
		map[string]interface{}{"foo": "bar"},
		false,
		csjwt.ValidationErrorSignatureInvalid,
		&csjwt.Parser{ValidMethods: []string{"HS256"}},
	},
	{
		"valid signing method",
		nil,
		defaultKeyFunc,
		map[string]interface{}{"foo": "bar"},
		true,
		0,
		&csjwt.Parser{ValidMethods: []string{"RS256", "HS256"}},
	},
}

func makeSample(c csjwt.MapClaims) []byte {
	key := csjwt.WithRSAPrivateKeyFromFile("test/sample_key")
	token := csjwt.NewWithClaims(csjwt.SigningMethodRS256, c)
	s, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}
	return s
}

func TestParser_Parse(t *testing.T) {
	for _, data := range jwtTestData {
		if len(data.tokenString) == 0 {
			data.tokenString = makeSample(data.claims)
		}

		var token csjwt.Token
		var err error
		if data.parser != nil {
			token, err = data.parser.Parse(data.tokenString, data.keyfunc)
		} else {
			token, err = csjwt.Parse(data.tokenString, data.keyfunc)
		}

		if !reflect.DeepEqual(&data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}
		if data.errors != 0 {
			if err == nil {
				t.Errorf("[%v] Expecting error.  Didn't get one.", data.name)
			} else {
				// compare the bitfield part of the error
				if e := err.(*csjwt.ValidationError).Errors; e != data.errors {
					t.Errorf("[%v] Errors don't match expectation.  %v != %v", data.name, e, data.errors)
				}
			}
		}
		if data.valid && len(token.Signature) == 0 {
			t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
		}
	}
}

func TestParseRequest(t *testing.T) {
	// Bearer token request
	for _, data := range jwtTestData {
		// FIXME: custom parsers are not supported by this helper.  skip tests that require them
		if data.parser != nil {
			t.Logf("Skipping [%v].  Custom parsers are not supported by ParseRequest", data.name)
			continue
		}

		if len(data.tokenString) == 0 {
			data.tokenString = makeSample(data.claims)
		}

		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.tokenString))
		token, err := csjwt.ParseFromRequestWithClaims(r, data.keyfunc, &csjwt.MapClaims{})

		if token.Raw == nil {
			t.Errorf("[%v] Token was not found: %v", data.name, err)
			continue
		}
		if !reflect.DeepEqual(&data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}
		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %v", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}
	}
}

func TestParseFromRequest(t *testing.T) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS512!`))
	clm := csjwt.MapClaims{
		"foo":               "bar",
		"user_id":           "hello_gophers",
		"cart_items":        "234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;",
		"last_viewed_items": "234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;",
		"requested_price":   float64(3.141592 * 2.718281 / 3),
		"checkout_step":     float64(3),
		"payment_valid":     true,
		"exp":               float64(time.Now().Add(time.Hour * 72).Unix()),
	}
	token := csjwt.NewWithClaims(csjwt.SigningMethodHS512, clm)
	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	rToken, err := csjwt.ParseFromRequest(r, func(t csjwt.Token) (csjwt.Key, error) {
		if have, want := t.Method.Alg(), csjwt.SigningMethodHS512.Alg(); have != want {
			return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
		}
		return key, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, rToken.Claims, &clm)
}

func benchmarkSigning(b *testing.B, method csjwt.Signer, key csjwt.Key) {
	t := csjwt.New(method)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := t.SignedString(key); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkMethodVerify(b *testing.B, method csjwt.Signer, signingString []byte, signature []byte, key csjwt.Key) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := method.Verify(signingString, signature, key); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkParseFromRequest_HS256(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS256!`))
	hmacFast, err := csjwt.NewHMACFast256(key)
	if err != nil {
		b.Fatal(err)
	}
	csjwt.RegisterSigningMethod(hmacFast)
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS256,
		key,
		func(t csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Method.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}
func BenchmarkParseFromRequest_HS384(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS384!`))
	hmacFast, err := csjwt.NewHMACFast384(key)
	if err != nil {
		b.Fatal(err)
	}
	csjwt.RegisterSigningMethod(hmacFast)
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS384,
		key,
		func(t csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Method.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}
func BenchmarkParseFromRequest_HS512(b *testing.B) {
	key := csjwt.WithPassword([]byte(`csjwt.SigningMethodHS512!`))
	hmacFast, err := csjwt.NewHMACFast512(key)
	if err != nil {
		b.Fatal(err)
	}
	csjwt.RegisterSigningMethod(hmacFast)
	benchmarkParseFromRequest(
		b,
		hmacFast, // csjwt.SigningMethodHS512,
		key,
		func(t csjwt.Token) (csjwt.Key, error) {
			if have, want := t.Method.Alg(), hmacFast.Alg(); have != want {
				return csjwt.Key{}, fmt.Errorf("Have: %s Want: %s", have, want)
			}
			return key, nil
		},
	)
}

func benchmarkParseFromRequest(b *testing.B, sm csjwt.Signer, key csjwt.Key, keyFunc csjwt.Keyfunc) {
	clm := csjwt.MapClaims{
		"foo":               "bar",
		"user_id":           "hello_gophers",
		"cart_items":        "234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;",
		"last_viewed_items": "234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;234234,12;34234,34;234234,1;123123,12;",
		"requested_price":   3.141592 * 2.718281 / 3,
		"checkout_step":     3,
		"payment_valid":     true,
		"exp":               time.Now().Add(time.Hour * 72).Unix(),
	}
	token := csjwt.NewWithClaims(sm, clm)
	tokenString, err := token.SignedString(key)
	if err != nil {
		b.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rToken, err := csjwt.ParseFromRequest(r, keyFunc)
		if err != nil {
			b.Fatal(err)
		}
		if !rToken.Valid {
			b.Fatalf("Token not valid: %#v", rToken)
		}
	}
	//b.Log("GC Pause:", gcPause())
}

func gcPause() time.Duration {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	return stats.Pause[0]
}
