package api

import (
	"testing"

	"github.com/Unleash/unleash-client-go/v3/context"
	"github.com/stretchr/testify/suite"
)

type VariantTestSuite struct {
	suite.Suite
	VariantWithOverride []Variant
	VariantWithoutOverride []Variant
}

func (suite *VariantTestSuite) SetupTest() {
	suite.VariantWithOverride = []Variant{
		Variant{
			Name:   "VarA",
			Weight: 33,
			Payload: Payload{
				Type:  "string",
				Value: "Test 1",
			},
			Overrides: []Override{
				Override{
					ContextName: "userId",
					Values: []string{
						"1",
					},
				},
				Override{
					ContextName: "sessionId",
					Values: []string{
						"ABCDE",
					},
				},
			},
		},
		Variant{
			Name:   "VarB",
			Weight: 33,
			Payload: Payload{
				Type:  "string",
				Value: "Test 2",
			},
			Overrides: []Override{
				Override{
					ContextName: "remoteAddress",
					Values: []string{
						"127.0.0.1",
					},
				},
			},
		},
		Variant{
			Name:   "VarC",
			Weight: 34,
			Payload: Payload{
				Type:  "string",
				Value: "Test 3",
			},
			Overrides: []Override{
				Override{
					ContextName: "env",
					Values: []string{
						"dev",
					},
				},
			},
		},
	}

	suite.VariantWithoutOverride = []Variant{
		Variant{
			Name:   "VarD",
			Weight: 33,
		},
		Variant{
			Name:   "VarE",
			Weight: 33,
		},
		Variant{
			Name:   "VarF",
			Weight: 34,
		},
	}	
}

func (suite *VariantTestSuite) TestGetVariantWhenFeatureHasNoVariant() {
	mockFeature := Feature{
		Name:    "test.variants",
		Enabled: true,
	}
	mockContext := &context.Context{}
	suite.Equal(DISABLED_VARIANT, mockFeature.GetVariant(mockContext), "Should return default variant")
}

func (suite *VariantTestSuite) TestGetVariantWhenFeatureIsNotEnabled() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  false,
		Variants: suite.VariantWithOverride,
	}
	mockContext := &context.Context{
		UserId:        "1",
		SessionId:     "ABCDE",
		RemoteAddress: "127.0.0.1",
	}
	suite.Equal(DISABLED_VARIANT, mockFeature.GetVariant(mockContext), "Should return default variant")
}

func (suite *VariantTestSuite) TestGetVariant_OverrideOnUserId() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithOverride,
	}
	mockContext := &context.Context{
		UserId:        "1",
		SessionId:     "ABCDE",
		RemoteAddress: "127.0.0.1",
	}
	expectedPayload := Payload{
		Type:  "string",
		Value: "Test 1",
	}
	suite.Equal("VarA", mockFeature.GetVariant(mockContext).Name, "Should return VarA")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
	suite.Equal(expectedPayload, mockFeature.GetVariant(mockContext).Payload, "Should be equal")
}


func (suite *VariantTestSuite) TestGetVariant_OverrideOnRemoteAddress() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithOverride,
	}
	mockContext := &context.Context{
		SessionId:     "FGHIJ",
		RemoteAddress: "127.0.0.1",
	}
	expectedPayload := Payload{
		Type:  "string",
		Value: "Test 2",
	}
	suite.Equal("VarB", mockFeature.GetVariant(mockContext).Name, "Should return VarB")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
	suite.Equal(expectedPayload, mockFeature.GetVariant(mockContext).Payload, "Should be equal")
}

func (suite *VariantTestSuite) TestGetVariant_OverrideOnSessionId() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithOverride,
	}
	mockContext := &context.Context{
		UserId:        "123",
		SessionId:     "ABCDE",
		RemoteAddress: "127.0.0.1",
	}
	expectedPayload := Payload{
		Type:  "string",
		Value: "Test 1",
	}
	suite.Equal("VarA", mockFeature.GetVariant(mockContext).Name, "Should return VarA")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
	suite.Equal(expectedPayload, mockFeature.GetVariant(mockContext).Payload, "Should be equal")
}

func (suite *VariantTestSuite) TestGetVariant_OverrideOnCustomProperties() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithOverride,
	}
	mockContext := &context.Context{
		Properties: map[string]string{
			"env": "dev",
		},
	}
	expectedPayload := Payload{
		Type:  "string",
		Value: "Test 3",
	}
	suite.Equal("VarC", mockFeature.GetVariant(mockContext).Name, "Should return VarC")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
	suite.Equal(expectedPayload, mockFeature.GetVariant(mockContext).Payload, "Should be equal")
}

func (suite *VariantTestSuite) TestGetVariant_ShouldReturnVarD() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithoutOverride,
	}
	mockContext := &context.Context{
		UserId:        "40",
	}
	suite.Equal("VarD", mockFeature.GetVariant(mockContext).Name, "Should return VarD")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
}

func (suite *VariantTestSuite) TestGetVariant_ShouldReturnVarE() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithoutOverride,
	}
	mockContext := &context.Context{
		UserId:        "123",
	}
	suite.Equal("VarE", mockFeature.GetVariant(mockContext).Name, "Should return VarE")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
}

func (suite *VariantTestSuite) TestGetVariant_ShouldReturnVarF() {
	mockFeature := Feature{
		Name:     "test.variants",
		Enabled:  true,
		Variants: suite.VariantWithoutOverride,
	}
	mockContext := &context.Context{
		UserId:        "163",
	}
	suite.Equal("VarF", mockFeature.GetVariant(mockContext).Name, "Should return VarF")
	suite.Equal(true, mockFeature.GetVariant(mockContext).Enabled, "Should be equal")
}


func TestVariantSuite(t *testing.T) {
	ts := VariantTestSuite{}
	suite.Run(t, &ts)
}
