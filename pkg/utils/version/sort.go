package version

// ByVersion implements sort.Interface based on the version information.
type ByVersion []string

func (versions ByVersion) Len() int {
	return len(versions)
}

func (versions ByVersion) Less(i, j int) bool {
	return CompareVersions(versions[i], versions[j]) < 1
}

func (versions ByVersion) Swap(i, j int) {
	versions[i], versions[j] = versions[j], versions[i]
}
