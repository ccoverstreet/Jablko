// Jablko Version 0.3.0
// Cale Overstreet
// Jun. 13, 2021

// Helper file for developing Jablko Applications. This file contains
// standardized functions for use in JMODs.

package main

import (
	"bytes"
	"net/http"
)

func JablkoSaveConfig(corePort string, modPort string, modKey string, config []byte) error {
	client := &http.Client{}
	reqPath := "http://localhost:" + corePort + "/service/saveConfig"

	req, err := http.NewRequest("POST", reqPath, bytes.NewBuffer(config))
	if err != nil {
		return err
	}

	req.Header.Add("JMOD-KEY", modKey)
	req.Header.Add("JMOD-PORT", modPort)

	_, err = client.Do(req)

	return err
}
