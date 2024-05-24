package models

import (
	"strings"
)

type RouteProp struct {
	Path    string
	Headers map[string]string
	Method  string
}

type RouteAction struct {
	RouteTo string `mapstructure:"route_to" json:"route_to" yaml:"route_to"`
}

type RouteCondition struct {
	PathPrefix string            `mapstructure:"path_prefix" json:"path_prefix" yaml:"path_prefix"`
	Headers    map[string]string `mapstructure:"headers" json:"headers" yaml:"headers"`
	Method     string            `mapstructure:"method" json:"method" yaml:"method"`
}

type RoutingRule struct {
	Conditions []RouteCondition `mapstructure:"conditions" json:"conditions" yaml:"conditions"`
	Actions    RouteAction      `mapstructure:"actions" json:"actions" yaml:"actionss"`
}
type Routing struct {
	Rules []RoutingRule `mapstructure:"rules" json:"rules" yaml:"rules"`
}

func (routeCondition *RouteCondition) Match(req *RouteProp) bool {
	if routeCondition.PathPrefix != "" && !strings.HasPrefix(req.Path, routeCondition.PathPrefix) {
		return false
	}

	if routeCondition.Method != "" && !strings.EqualFold(req.Method, routeCondition.Method) {
		return false
	}

	for k, v := range routeCondition.Headers {
		if value, ok := req.Headers[k]; ok {
			if v != value {
				return false
			}
		}
	}

	return true
}

func (rule *RoutingRule) Match(req *RouteProp) bool {
	for _, condition := range rule.Conditions {
		if !condition.Match(req) {
			return false
		}
	}

	return true
}

func (routing *Routing) GetRoute(req *RouteProp) string {
	for _, rule := range routing.Rules {
		if rule.Match(req) {
			return rule.Actions.RouteTo
		}
	}

	return ""
}
