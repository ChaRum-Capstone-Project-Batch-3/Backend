package report
func (ru *ReportUseCase) Create(domain *Domain) (Domain, error) {
	// check ReportedID if exist in users or threads ID
	reportType, err := ru.CheckID(domain.ReportedID)
	if err != nil {
		return Domain{}, errors.New("ID not found")
	}
	_, err = ru.reportRepository.CheckByUserID(domain.UserID, domain.ReportedID)
	if err == nil {
		return Domain{}, errors.New("already reported")
	}

	domain.Id = primitive.NewObjectID()
	domain.ReportType = reportType
	domain.ReportDetail = "Inappropriate content or behavior"
	domain.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	domain.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	report, err := ru.reportRepository.Create(domain)
	if err != nil {
		return Domain{}, err
	}

	return report, nil
}
