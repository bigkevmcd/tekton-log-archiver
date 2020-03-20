package archiver

// Archivable returns true if this run should be archived.
func Archivable(ag annotationsGetter) bool {
	for k, v := range ag.Annotations() {
		if k == ArchivableName && v == "true" {
			return true
		}
	}
	return false
}
