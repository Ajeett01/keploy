package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"go.keploy.io/server/pkg/service/serve/graph/model"
	"go.uber.org/zap"
)

// RunTestSet is the resolver for the runTestSet field.
func (r *mutationResolver) RunTestSet(ctx context.Context, testSet string) (*model.RunTestSetResponse, error) {
	if r.Resolver == nil {
		err := fmt.Errorf(Emoji + "failed to get Resolver")
		return nil, err
	}
	
	tester := r.Resolver.Tester

	if tester == nil {
		r.Logger.Error("failed to get tester from resolver")
		return nil, fmt.Errorf(Emoji+"failed to run testSet:%v", testSet)
	}

	testRunChan := make(chan string, 1)
	pid := r.Resolver.AppPid
	testCasePath := r.Resolver.Path
	testReportPath := r.Resolver.TestReportPath
	delay := r.Resolver.Delay

	// since we are not restarting the application again and again,
	// so we don't want for the test to wait for the delay time everytime a new test runs.
	if r.firstRequestDone {
		delay = 1
	}

	if !r.Resolver.firstRequestDone {
		r.firstRequestDone = true
	}

	testReportFS := r.Resolver.TestReportFS
	if tester == nil {
		r.Logger.Error( "failed to get testReportFS from resolver")
		return nil, fmt.Errorf(Emoji+"failed to run testSet:%v", testSet)
	}

	ys := r.Resolver.YS
	if ys == nil {
		r.Logger.Error( "failed to get ys from resolver")
		return nil, fmt.Errorf(Emoji+"failed to run testSet:%v", testSet)
	}

	loadedHooks := r.LoadedHooks
	if loadedHooks == nil {
		r.Logger.Error( "failed to get loadedHooks from resolver")
		return nil, fmt.Errorf(Emoji+"failed to run testSet:%v", testSet)
	}

	go func() {
		r.Logger.Debug("starting testrun...", zap.Any("testSet", testSet))
		tester.RunTestSet(testSet, testCasePath, testReportPath, "", "", "", delay, pid, ys, loadedHooks, testReportFS, testRunChan)
	}()

	testRunID := <-testRunChan
	r.Logger.Debug("", zap.Any("testRunID", testRunID))

	return &model.RunTestSetResponse{Success: true, TestRunID: testRunID}, nil
}

// TestSets is the resolver for the testSets field.
func (r *queryResolver) TestSets(ctx context.Context) ([]string, error) {
	if r.Resolver == nil {
		err := fmt.Errorf(Emoji + "failed to get Resolver")
		return nil, err
	}
	testPath := r.Resolver.Path

	testSets, err := r.Resolver.YS.ReadSessionIndices(testPath)
	if err != nil {
		r.Resolver.Logger.Error("failed to fetch test sets", zap.Any("testPath", testPath), zap.Error(err))
		return nil, err
	}

	// Print debug log for retrieved qualified test sets
	if len(testSets) > 0 {
		r.Resolver.Logger.Debug(fmt.Sprintf("Retrieved test sets: %v", testSets), zap.Any("testPath", testPath))
	} else {
		r.Resolver.Logger.Debug("No test sets found", zap.Any("testPath", testPath))
	}

	return testSets, nil
}

// TestSetStatus is the resolver for the testSetStatus field.
func (r *queryResolver) TestSetStatus(ctx context.Context, testRunID string) (*model.TestSetStatus, error) {
	if r.Resolver == nil {
		err := fmt.Errorf(Emoji + "failed to get Resolver")
		return nil, err
	}
	testReportFs := r.Resolver.TestReportFS

	if testReportFs == nil {
		r.Logger.Error( "failed to get testReportFS from resolver")
		return nil, fmt.Errorf(Emoji+"failed to get the status for testRunID:%v", testRunID)
	}
	testReport, err := testReportFs.Read(ctx, r.Resolver.TestReportPath, testRunID)
	if err != nil {
		r.Logger.Error("failed to fetch testReport", zap.Any("testRunID", testRunID), zap.Error(err))
		return nil, err
	}

	r.Logger.Debug("", zap.Any("testRunID", testRunID), zap.Any("testSetStatus", testReport.Status))
	return &model.TestSetStatus{Status: testReport.Status}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }