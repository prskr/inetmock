package main

func SnapshotBuild() error {
	return GoReleaser("release", "--snapshot", "--skip-publish", "--rm-dist")
}
