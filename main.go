package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
)

func main() {
	if err := test(); err != nil {
		log.Fatal(err)
	}
}

func test() error {
	provider, err := auth.InstancePrincipalConfigurationProviderWithCustomClient(
		func(dispatcher common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error) {
			return &http.Client{
				Timeout: 5 * time.Second,
			}, nil
		})
	if err != nil {
		return fmt.Errorf("making provider: %w", err)
	}

	req := common.MakeDefaultHTTPRequest(http.MethodGet, "v1/instancePrincipalRootCACertificates")
	req.URL.Scheme = "https"
	req.URL.Host = "auth.us-phoenix-1.oraclecloud.com"

	signedHeaders := common.DefaultGenericHeaders()
	signer := common.RequestSigner(provider, signedHeaders, common.DefaultBodyHeaders())
	signer.Sign(&req)

	dump, err := httputil.DumpRequestOut(&req, true)
	if err != nil {
		return fmt.Errorf("dumping request: %w", err)
	}
	fmt.Println(string(dump))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(&req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	fmt.Println("response Body:", string(body))

	return nil
}
