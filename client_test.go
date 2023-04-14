package gofantasy

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

var accessToken = "JKjQl1eYtF31dtI46eoN.fSA9vCsBz14KYU2MNooM2pWoylgjBug2O9D0W8zc4DEsEkNINoRlivoMApEgJGqjjC_v7ODBjgjJYbPZ33QOhr4e85UqJvM5Oam3pViN17GYn0anaJpXN0Oz6f6DmQrazlvU5k3K4XXZCgbml.V0HZ1vCGgcltm843x2VvD9J.X7ZUgbl1GYPGaPceyQHBGwW2DdyG3EKNKLrYjiV9I.lsojHsmF3IPVkXanqSaE3Ui_rHeACthCEJBid4weoXd5k2Av2GK0lradeR3blSAJO0u4Keb734_72Z8qOpapjwj2yxg8J9aNTgMVoGWg6wLxn_BBVvvSpJZvz0qSifgD8BqgXHZ.ZE1THdeSV97MqknrBZXyO83ltidAbi3N5u.R5spwRDXZ_Y8j4pD3JOE1mQn412gHdAW24ibjVu6IBHU9vrMVIE9yWSk6u3P0ATFLdhT6Rxq1D5hOf71r1p7QorZkb6FCvym_lLbdjoPj3fXL6bfvAXielk.GWwbpIxDj7RTHZv8ERYlK.bhPkfXq5Kvxe8pQkitLWh8jmH4YNVsge9QHCYZHkfnWtE77FtdbWOY66cFVxRBpsk_9JFtyf5gqkDV5VrH7EpGCf3A8PyAgS0U6XfjdsY97V..Ecq2tCYsFHG4oIR7DPe_RZePVoI6hMgsO9oz53Jbe3wUy6r9LRTfp34DgM1aLSLwzGxdw6oC8TO1Iypx7ie50iUlS0txixYB5Rn1G44L2M7eT.HABB7oXcMJur_Ah1zRuprsgGZ79np5q64Updp4OLh2gjS6qMm2sLEgjxKewb2D6y3oT0I7kFVmm7K0RVr3e059Kqq0x3tLHolyAInD6oKCMsZ55YtyeyHxoKyGLp.fA8jQciJ_XItmBNLNH.o0_ntjpoch99LRtfXI7zXSjrg7MyQg4j11CXAHuMycgJaR9t1asH92XEPeDL8dnumNQHckYkHkTZxN_LbaFA.jO9bF6eOPtO252Wlp.tJeKyGJ8DW0jQqGBtUISamaXNvvLVVSsex2LlrHUTZ_"

func Test_SendRequest(t *testing.T) {
	c := NewClient(accessToken)
	ctx := context.Background()
	fmt.Println(c.fullURL(""))
	req, err := c.requestBuilder.build(ctx, http.MethodGet, c.fullURL(""), nil)
	if err != nil {
		return
	}
	response := &FantasyContent{}
	err = c.sendRequest(req, response)
	if err != nil {
		return
	}
	fmt.Println("-------", response)
}
