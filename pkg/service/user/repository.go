package user
//
//import (
//	"context"
//	"github.com/opentracing/opentracing-go"
//	"github.com/openzipkin/zipkin-go"
//	"github.com/openzipkin/zipkin-go/model"
//	"gitlab.com/vdat/mcsvc/sidm/pkg/idmservice"
//)
//
//type Repository interface {
//	GetUserByUsername(ctx context.Context, Username string) (*User, error)
//}
//
//
///**
// * In-memory User Repository
// * useful for development
// */
//type inmemRepository struct {
//	users []User
//	// tracers (optional)
//	tracer opentracing.Tracer
//	zipkinTracer *zipkin.Tracer
//}
//
//// GetUserByUsername - return a single User pointer, nil if not found
//func (i inmemRepository) GetUserByUsername(ctx context.Context, Username string) (*User, error) {
//	// setup tracing (for monitor purpose)
//	if i.zipkinTracer != nil {
//		var sc model.SpanContext
//		if parentSpan := zipkin.SpanFromContext(ctx); parentSpan != nil {
//			sc = parentSpan.Context()
//		}
//		sp := i.zipkinTracer.StartSpan("GetUserByUsername", zipkin.Parent(sc))
//		defer sp.Finish()
//
//		ctx = zipkin.NewContext(ctx, sp)
//	}
//
//	// actual code
//	var result User
//
//	for _, u := range i.users {
//		if u.Username == Username {
//			result = u
//			break
//		}
//	}
//
//	return &result, nil
//}
//
//func NewInmemUserRepository(tracer opentracing.Tracer, zipkinTracer *zipkin.Tracer) Repository {
//	return inmemRepository{
//		[]User{
//			{ID: "1", Username: "user1", First: "User1"},
//			{ID: "2", Username: "user2", First: "User2"},
//		},
//		tracer,
//		zipkinTracer,
//	}
//}
//
///**
// * IdM Repository implementation
// *
// * info: IdM (Identity Management Service) is a service manage user information in VDAT's ecosystem
// */
//type idmRepository struct {
//	client *idmservice.Service
//	// tracers
//	tracer opentracing.Tracer
//	zipkinTracer *zipkin.Tracer
//}
//
//// GetUserByUsername - return a single User pointer, nil if not found
//func (i idmRepository) GetUserByUsername(ctx context.Context, Username string) (*User, error)  {
//	// setup tracing
//	if i.zipkinTracer != nil {
//		var sc model.SpanContext
//		if parentSpan := zipkin.SpanFromContext(ctx); parentSpan != nil {
//			sc = parentSpan.Context()
//		}
//		sp := i.zipkinTracer.StartSpan("GetUserByUsername", zipkin.Parent(sc))
//		defer sp.Finish()
//
//		ctx = zipkin.NewContext(ctx, sp)
//	}
//
//	panic("implement me")
//}
//
//func NewIdmUserRepository(idm *idmservice.Service, tracer opentracing.Tracer, zipkinTracer *zipkin.Tracer) Repository {
//	return idmRepository{
//		idm,
//		tracer,
//		zipkinTracer,
//	}
//}
