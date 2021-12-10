package clientapi

/*
func TestRaHub(t *testing.T) {
	test.CreateServerConfigFile()
	// We can't use this default config, because the go test
	// will run TestRaHub and TestClientAPI parallelly, so
	// it can't bind the same port on the same host.
	// Here we just use another port to bind.
	//cfg := config.GetDefault()
	//server := cfg.GetPort()
	server := ":40004"
	defer test.RemoveConfigFile()
	vm, err := verifier.CreateVerifierMgr()
	if err != nil {
		fmt.Println(err)
		return
	}
	cm := cache.CreateCacheMgr(cache.DEFAULTRACNUM, vm)
	go StartServer(server, cm)

	const addrRaHub string = ":40003"
	go StartRaHub(addrRaHub, server)

	_, err = DoGenerateIKCert(addrRaHub, &GenerateIKCertRequest{
		EkCert: certPEM,
		IkPub:  pubPEM,
		IkName: nil,
	})
	if err != nil {
		t.Errorf("Client: invoke CreateIKCert error %v", err)
	}
	t.Logf("Client: invoke CreateIKCert ok")

	ci, err := json.Marshal(map[string]string{"test name": "test value"})
	if err != nil {
		t.Error(err)
	}
	r, err := DoRegisterClient(addrRaHub, &RegisterClientRequest{
		Ic:         &Cert{Cert: createRandomCert()},
		ClientInfo: &ClientInfo{ClientInfo: string(ci)},
	})
	if err != nil {
		t.Errorf("Client: invoke RegisterClient error %v", err)
	}
	t.Logf("Client: invoke RegisterClient ok, clientID=%d", r.GetClientId())

	_, err = DoSendHeartbeat(addrRaHub, &SendHeartbeatRequest{ClientId: r.GetClientId()})
	if err != nil {
		t.Errorf("Client: invoke SendHeartbeat error %v", err)
	}
	t.Logf("Client: invoke SendHeartbeat ok")

	cfg := config.GetDefault(config.ConfServer)
	trustmgr.SetValidator(&testValidator{})
	_, err = DoSendReport(addrRaHub, &SendReportRequest{ClientId: r.GetClientId(),
		TrustReport: &TrustReport{
			PcrInfo: &PcrInfo{PcrValues: map[int32]string{
				1: "pcr value1",
				2: "pcr value2",
			},
				PcrQuote: &PcrQuote{
					Quoted: []byte("test quote"),
				},
				Algorithm: cfg.GetDigestAlgorithm(),
			},
			Manifest: []*Manifest{},
			ClientId: r.GetClientId(),
		}})
	if err != nil {
		t.Errorf("Client: invoke SendReport error %v", err)
	}
	t.Logf("Client: invoke SendReport ok")

	u, err := DoUnregisterClient(addrRaHub,
		&UnregisterClientRequest{ClientId: r.GetClientId()})
	if err != nil {
		t.Errorf("Client: invoke UnregisterClient error %v", err)
	}
	t.Logf("Client: invoke UnregisterClient %v", u.Result)

}
*/
