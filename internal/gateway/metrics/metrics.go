package metrics

const prometheusNamespace = "dostavkee"

const (
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

func RegisterMetrics() error {
	if err := registerOrderMetrics(); err != nil {
		return err
	}

	if err := registerExternalOrderMetrics(); err != nil {
		return err
	}

	if err := registerPaymentMetrics(); err != nil {
		return err
	}

	if err := registerDgisMetrics(); err != nil {
		return err
	}

	if err := registerDadataMetrics(); err != nil {
		return err
	}

	if err := registerDriveeMetrics(); err != nil {
		return err
	}

	if err := registerFileStorageMetrics(); err != nil {
		return err
	}

	if err := registerSmtpSendMetrics(); err != nil {
		return err
	}

	if err := registerTelegramMetrics(); err != nil {
		return err
	}

	if err := registerCompanyMetrics(); err != nil {
		return err
	}

	if err := registerInvoiceMetrics(); err != nil {
		return err
	}

	if err := registerTinkoffMetrics(); err != nil {
		return err
	}

	if err := registerCorpInvoiceMetrics(); err != nil {
		return err
	}

	if err := registerGotenbergMetrics(); err != nil {
		return err
	}

	return nil
}
