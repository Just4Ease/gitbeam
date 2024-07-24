//go:generate mockgen -source=api/pb/commits/commits.pb.go -destination=mocks/commits_monitor_service_mock.go -package=mocks
//go:generate mockgen -source=api/pb/repos/repos.pb.go -destination=mocks/repos_manager_service_mock.go -package=mocks

package main
