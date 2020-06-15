package health

func StaticResultCheck(status Status) Check {
	return func() CheckResult {
		return CheckResult{
			Status: status,
		}
	}
}

func StaticResultCheckWithMessage(status Status, message string) Check {
	return func() CheckResult {
		return CheckResult{
			Status:  status,
			Message: message,
		}
	}
}
