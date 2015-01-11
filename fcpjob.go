package fcp

type fcpJob struct {
}

type fileFetchJob struct {
	fcpJob
}

type fileInsertJob struct {
	fcpJob
}

type siteInsertJob struct {
	fileInsertJob
}

type siteUpdateJob struct {
	siteInsertJob
}
