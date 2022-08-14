package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	entryPointXSD := "../../_random_stuff/xbrl/us-gaap-2022/us-gaap-2022/entire/us-gaap-entryPoint-all-2022.xsd"
	entryPointRootDir := path.Dir(entryPointXSD)

	reader, err := os.Open(entryPointXSD)
	require.NoError(t, err)

	xmlBytes, err := io.ReadAll(reader)
	require.NoError(t, err)

	var taxonomy XBRLTaxonomy
	require.NoError(t, xml.Unmarshal(xmlBytes, &taxonomy))

	for _, elem := range taxonomy.Elements {
		elem.OriginalFile = entryPointXSD
	}

	for _, linkbase := range taxonomy.AppInfo.LinkBases {
		linkbase.OriginalFile = entryPointXSD
	}

	require.NoError(t, taxonomy.ResolveLocalImports(entryPointRootDir, true))

	t.Logf("total elements: %d", len(taxonomy.Elements))
	t.Logf("files imported: %d", len(ResolvedImportsSet))

	for _, elem := range taxonomy.Elements {
		if elem.Name == "Assets" {
			t.Log(elem)
		}
	}

	//require.NoError(t, taxonomy.ResolveLocalLinkBases(entryPointRootDir, true))

	t.Logf("total linkbases: %d", len(taxonomy.AppInfo.LinkBases))
	t.Logf("linkbases imported: %d", len(ResolvedLinkBaseSet))

	start := time.Now()
	totalLocators := 0

	labels, err := taxonomy.LookupLabels("Assets")
	require.NoError(t, err)

	for _, label := range labels {
		t.Logf("label for Assets: %s (%s)", label.Value, label.Role)
	}

	t.Logf("finding labels took %s", time.Since(start).String())
	t.Logf("total of %d locators", totalLocators)

	require.NoError(t, FormatLegendForPresentationLink(&taxonomy, "http://fasb.org/us-gaap/role/statement/StatementOfFinancialPositionClassified"))
}

func TestFormatThing(t *testing.T) {
	presLink := `
<link:presentationLink xlink:role='http://fasb.org/us-gaap/role/statement/StatementOfFinancialPositionClassified' xlink:type='extended'>
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsPayableAndAccruedLiabilitiesCurrentAbstract' xlink:label='loc_AccountsPayableAndAccruedLiabilitiesCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsPayableCurrent' xlink:label='loc_AccountsPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedLiabilitiesCurrentAbstract' xlink:label='loc_AccruedLiabilitiesCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsPayableAndAccruedLiabilitiesCurrent' xlink:label='loc_AccountsPayableAndAccruedLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_AccountsReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_AccountsReceivableExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_AccountsReceivableAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_AccountsReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableNetCurrentAbstract' xlink:label='loc_AccountsReceivableNetCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableGrossCurrent' xlink:label='loc_AccountsReceivableGrossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AllowanceForDoubtfulAccountsReceivableCurrent' xlink:label='loc_AllowanceForDoubtfulAccountsReceivableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableNetCurrent' xlink:label='loc_AccountsReceivableNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableNetNoncurrentAbstract' xlink:label='loc_AccountsReceivableNetNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableGrossNoncurrent' xlink:label='loc_AccountsReceivableGrossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AllowanceForDoubtfulAccountsReceivableNoncurrent' xlink:label='loc_AllowanceForDoubtfulAccountsReceivableNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableNetNoncurrent' xlink:label='loc_AccountsReceivableNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsNotesAndLoansReceivableNetCurrentAbstract' xlink:label='loc_AccountsNotesAndLoansReceivableNetCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansReceivableNetCurrentAbstract' xlink:label='loc_NotesAndLoansReceivableNetCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsNotesAndLoansReceivableNetCurrent' xlink:label='loc_AccountsNotesAndLoansReceivableNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermAccountsNotesAndLoansReceivableNetNoncurrentAbstract' xlink:label='loc_LongTermAccountsNotesAndLoansReceivableNetNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansReceivableNetNoncurrentAbstract' xlink:label='loc_NotesAndLoansReceivableNetNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermAccountsNotesAndLoansReceivableNetNoncurrent' xlink:label='loc_LongTermAccountsNotesAndLoansReceivableNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EmployeeRelatedLiabilitiesCurrent' xlink:label='loc_EmployeeRelatedLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TaxesPayableCurrentAbstract' xlink:label='loc_TaxesPayableCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InterestAndDividendsPayableCurrent' xlink:label='loc_InterestAndDividendsPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedLiabilitiesCurrent' xlink:label='loc_AccruedLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedOtherComprehensiveIncomeLossNetOfTaxAbstract' xlink:label='loc_AccumulatedOtherComprehensiveIncomeLossNetOfTaxAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedOtherComprehensiveIncomeLossForeignCurrencyTranslationAdjustmentNetOfTax' xlink:label='loc_AccumulatedOtherComprehensiveIncomeLossForeignCurrencyTranslationAdjustmentNetOfTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedOtherComprehensiveIncomeLossAvailableForSaleSecuritiesAdjustmentNetOfTax' xlink:label='loc_AccumulatedOtherComprehensiveIncomeLossAvailableForSaleSecuritiesAdjustmentNetOfTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AociLossCashFlowHedgeCumulativeGainLossAfterTax' xlink:label='loc_AociLossCashFlowHedgeCumulativeGainLossAfterTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedOtherComprehensiveIncomeLossDefinedBenefitPensionAndOtherPostretirementPlansNetOfTax' xlink:label='loc_AccumulatedOtherComprehensiveIncomeLossDefinedBenefitPensionAndOtherPostretirementPlansNetOfTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedOtherComprehensiveIncomeLossFinancialLiabilityFairValueOptionAfterTax' xlink:label='loc_AccumulatedOtherComprehensiveIncomeLossFinancialLiabilityFairValueOptionAfterTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AociDerivativeQualifyingAsHedgeExcludedComponentAfterTax' xlink:label='loc_AociDerivativeQualifyingAsHedgeExcludedComponentAfterTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedOtherComprehensiveIncomeLossNetOfTax' xlink:label='loc_AccumulatedOtherComprehensiveIncomeLossNetOfTax' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdditionalPaidInCapitalAbstract' xlink:label='loc_AdditionalPaidInCapitalAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdditionalPaidInCapitalCommonStock' xlink:label='loc_AdditionalPaidInCapitalCommonStock' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdditionalPaidInCapitalPreferredStock' xlink:label='loc_AdditionalPaidInCapitalPreferredStock' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdditionalPaidInCapital' xlink:label='loc_AdditionalPaidInCapital' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetRetirementObligationsNoncurrentAbstract' xlink:label='loc_AssetRetirementObligationsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MineReclamationAndClosingLiabilityNoncurrent' xlink:label='loc_MineReclamationAndClosingLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OilAndGasReclamationLiabilityNoncurrent' xlink:label='loc_OilAndGasReclamationLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedCappingClosurePostClosureAndEnvironmentalCostsNoncurrent' xlink:label='loc_AccruedCappingClosurePostClosureAndEnvironmentalCostsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DecommissioningLiabilityNoncurrent' xlink:label='loc_DecommissioningLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SpentNuclearFuelObligationNoncurrent' xlink:label='loc_SpentNuclearFuelObligationNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetRetirementObligationsNoncurrent' xlink:label='loc_AssetRetirementObligationsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsHeldForSaleNotPartOfDisposalGroupCurrentAbstract' xlink:label='loc_AssetsHeldForSaleNotPartOfDisposalGroupCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TradeAndLoansReceivablesHeldForSaleNetNotPartOfDisposalGroup' xlink:label='loc_TradeAndLoansReceivablesHeldForSaleNetNotPartOfDisposalGroup' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsHeldForSaleNotPartOfDisposalGroupCurrentOther' xlink:label='loc_AssetsHeldForSaleNotPartOfDisposalGroupCurrentOther' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsHeldForSaleNotPartOfDisposalGroupCurrent' xlink:label='loc_AssetsHeldForSaleNotPartOfDisposalGroupCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsAbstract' xlink:label='loc_AssetsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsCurrentAbstract' xlink:label='loc_AssetsCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsNoncurrentAbstract' xlink:label='loc_AssetsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_Assets' xlink:label='loc_Assets' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsPledgingPurposeExtensibleEnumeration' xlink:label='loc_AssetsPledgingPurposeExtensibleEnumeration' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashCashEquivalentsAndShortTermInvestmentsAbstract' xlink:label='loc_CashCashEquivalentsAndShortTermInvestmentsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ReceivablesNetCurrentAbstract' xlink:label='loc_ReceivablesNetCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleExcludingAccruedInterestCurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryNetAbstract' xlink:label='loc_InventoryNetAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidExpenseCurrentAbstract' xlink:label='loc_PrepaidExpenseCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetNetCurrentAbstract' xlink:label='loc_ContractWithCustomerAssetNetCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CapitalizedContractCostNetCurrent' xlink:label='loc_CapitalizedContractCostNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherAssetsCurrent' xlink:label='loc_OtherAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCostsCurrentAbstract' xlink:label='loc_DeferredCostsCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesAbstract' xlink:label='loc_DerivativeInstrumentsAndHedgesAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsOfDisposalGroupIncludingDiscontinuedOperationCurrentAbstract' xlink:label='loc_AssetsOfDisposalGroupIncludingDiscontinuedOperationCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseCurrent' xlink:label='loc_NetInvestmentInLeaseCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RegulatoryAssetsCurrent' xlink:label='loc_RegulatoryAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FundsHeldForClients' xlink:label='loc_FundsHeldForClients' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRentAssetNetCurrent' xlink:label='loc_DeferredRentAssetNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdvancesOnInventoryPurchases' xlink:label='loc_AdvancesOnInventoryPurchases' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdvanceRoyaltiesCurrent' xlink:label='loc_AdvanceRoyaltiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DepositsAssetsCurrent' xlink:label='loc_DepositsAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BusinessCombinationContingentConsiderationAssetCurrent' xlink:label='loc_BusinessCombinationContingentConsiderationAssetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsCurrent' xlink:label='loc_AssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryNoncurrentAbstract' xlink:label='loc_InventoryNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OperatingLeaseRightOfUseAsset' xlink:label='loc_OperatingLeaseRightOfUseAsset' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinanceLeaseRightOfUseAsset' xlink:label='loc_FinanceLeaseRightOfUseAsset' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LeveragedLeasesNetInvestmentInLeveragedLeasesDisclosureInvestmentInLeveragedLeasesNet' xlink:label='loc_LeveragedLeasesNetInvestmentInLeveragedLeasesDisclosureInvestmentInLeveragedLeasesNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PropertyPlantAndEquipmentNetAbstract' xlink:label='loc_PropertyPlantAndEquipmentNetAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PropertyPlantAndEquipmentCollectionsNotCapitalized' xlink:label='loc_PropertyPlantAndEquipmentCollectionsNotCapitalized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OilAndGasPropertySuccessfulEffortMethodNet' xlink:label='loc_OilAndGasPropertySuccessfulEffortMethodNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OilAndGasPropertyFullCostMethodNet' xlink:label='loc_OilAndGasPropertyFullCostMethodNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermInvestmentsAndReceivablesNetAbstract' xlink:label='loc_LongTermInvestmentsAndReceivablesNetAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleExcludingAccruedInterestNoncurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseNoncurrent' xlink:label='loc_NetInvestmentInLeaseNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_Goodwill' xlink:label='loc_Goodwill' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_IntangibleAssetsNetExcludingGoodwillAbstract' xlink:label='loc_IntangibleAssetsNetExcludingGoodwillAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidExpenseNoncurrentAbstract' xlink:label='loc_PrepaidExpenseNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetNetNoncurrentAbstract' xlink:label='loc_ContractWithCustomerAssetNetNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CapitalizedContractCostNetNoncurrent' xlink:label='loc_CapitalizedContractCostNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherAssetsNoncurrent' xlink:label='loc_OtherAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesNoncurrentAbstract' xlink:label='loc_DerivativeInstrumentsAndHedgesNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RegulatedEntityOtherAssetsNoncurrentAbstract' xlink:label='loc_RegulatedEntityOtherAssetsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InvestmentsAndOtherNoncurrentAssets' xlink:label='loc_InvestmentsAndOtherNoncurrentAssets' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCostsNoncurrentAbstract' xlink:label='loc_DeferredCostsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DepositsAssetsNoncurrent' xlink:label='loc_DepositsAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InsuranceReceivableForMalpracticeNoncurrent' xlink:label='loc_InsuranceReceivableForMalpracticeNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredIncomeTaxAssetsNet' xlink:label='loc_DeferredIncomeTaxAssetsNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRentReceivablesNetNoncurrent' xlink:label='loc_DeferredRentReceivablesNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsHeldInTrustNoncurrent' xlink:label='loc_AssetsHeldInTrustNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DefinedBenefitPlanAssetsForPlanBenefitsNoncurrent' xlink:label='loc_DefinedBenefitPlanAssetsForPlanBenefitsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndInvestmentsNoncurrentAbstract' xlink:label='loc_RestrictedCashAndInvestmentsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsOfDisposalGroupIncludingDiscontinuedOperationNoncurrentAbstract' xlink:label='loc_AssetsOfDisposalGroupIncludingDiscontinuedOperationNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdvanceRoyaltiesNoncurrent' xlink:label='loc_AdvanceRoyaltiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BusinessCombinationContingentConsiderationAssetNoncurrent' xlink:label='loc_BusinessCombinationContingentConsiderationAssetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AmortizationMethodQualifiedAffordableHousingProjectInvestments' xlink:label='loc_AmortizationMethodQualifiedAffordableHousingProjectInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsNoncurrent' xlink:label='loc_AssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashAndCashEquivalentsAtCarryingValueAbstract' xlink:label='loc_CashAndCashEquivalentsAtCarryingValueAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_Cash' xlink:label='loc_Cash' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashEquivalentsAtCarryingValue' xlink:label='loc_CashEquivalentsAtCarryingValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashAndCashEquivalentsAtCarryingValue' xlink:label='loc_CashAndCashEquivalentsAtCarryingValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashAndCashEquivalentsPledgedStatusExtensibleEnumeration' xlink:label='loc_CashAndCashEquivalentsPledgedStatusExtensibleEnumeration' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashAndCashEquivalentsPledgingPurposeExtensibleEnumeration' xlink:label='loc_CashAndCashEquivalentsPledgingPurposeExtensibleEnumeration' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndInvestmentsCurrentAbstract' xlink:label='loc_RestrictedCashAndInvestmentsCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermInvestmentsAbstract' xlink:label='loc_ShortTermInvestmentsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CashCashEquivalentsAndShortTermInvestments' xlink:label='loc_CashCashEquivalentsAndShortTermInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StatementClassOfStockAxis' xlink:label='loc_StatementClassOfStockAxis' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ClassOfStockDomain' xlink:label='loc_ClassOfStockDomain' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonClassAMember' xlink:label='loc_CommonClassAMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonClassBMember' xlink:label='loc_CommonClassBMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonClassCMember' xlink:label='loc_CommonClassCMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CumulativePreferredStockMember' xlink:label='loc_CumulativePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NoncumulativePreferredStockMember' xlink:label='loc_NoncumulativePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemablePreferredStockMember' xlink:label='loc_RedeemablePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NonredeemablePreferredStockMember' xlink:label='loc_NonredeemablePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertiblePreferredStockMember' xlink:label='loc_ConvertiblePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredClassAMember' xlink:label='loc_PreferredClassAMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredClassBMember' xlink:label='loc_PreferredClassBMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesAPreferredStockMember' xlink:label='loc_SeriesAPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesBPreferredStockMember' xlink:label='loc_SeriesBPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesCPreferredStockMember' xlink:label='loc_SeriesCPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesDPreferredStockMember' xlink:label='loc_SeriesDPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesEPreferredStockMember' xlink:label='loc_SeriesEPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesFPreferredStockMember' xlink:label='loc_SeriesFPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesGPreferredStockMember' xlink:label='loc_SeriesGPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeriesHPreferredStockMember' xlink:label='loc_SeriesHPreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockNumberOfSharesParValueAndOtherDisclosuresAbstract' xlink:label='loc_CommonStockNumberOfSharesParValueAndOtherDisclosuresAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockParOrStatedValuePerShare' xlink:label='loc_CommonStockParOrStatedValuePerShare' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockNoParValue' xlink:label='loc_CommonStockNoParValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockSharesSubscriptions' xlink:label='loc_CommonStockSharesSubscriptions' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockSharesSubscribedButUnissued' xlink:label='loc_CommonStockSharesSubscribedButUnissued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockSharesAuthorized' xlink:label='loc_CommonStockSharesAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockSharesAuthorizedUnlimited' xlink:label='loc_CommonStockSharesAuthorizedUnlimited' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockSharesIssued' xlink:label='loc_CommonStockSharesIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockSharesOutstanding' xlink:label='loc_CommonStockSharesOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockValueOutstanding' xlink:label='loc_CommonStockValueOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockOtherSharesOutstanding' xlink:label='loc_CommonStockOtherSharesOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockOtherValueOutstanding' xlink:label='loc_CommonStockOtherValueOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockConversionBasis' xlink:label='loc_CommonStockConversionBasis' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockIssuedEmployeeStockTrust' xlink:label='loc_CommonStockIssuedEmployeeStockTrust' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockIssuedEmployeeTrustDeferred' xlink:label='loc_CommonStockIssuedEmployeeTrustDeferred' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockHeldInTrust' xlink:label='loc_CommonStockHeldInTrust' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetGrossCurrent' xlink:label='loc_ContractWithCustomerAssetGrossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetAccumulatedAllowanceForCreditLossCurrent' xlink:label='loc_ContractWithCustomerAssetAccumulatedAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetNetCurrent' xlink:label='loc_ContractWithCustomerAssetNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetGrossNoncurrent' xlink:label='loc_ContractWithCustomerAssetGrossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetAccumulatedAllowanceForCreditLossNoncurrent' xlink:label='loc_ContractWithCustomerAssetAccumulatedAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerAssetNetNoncurrent' xlink:label='loc_ContractWithCustomerAssetNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NonredeemableConvertiblePreferredStockMember' xlink:label='loc_NonredeemableConvertiblePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemableConvertiblePreferredStockMember' xlink:label='loc_RedeemableConvertiblePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContingentConvertiblePreferredStockMember' xlink:label='loc_ContingentConvertiblePreferredStockMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesAvailableForSaleAmortizedCostExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesCurrentAbstract' xlink:label='loc_DebtSecuritiesCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TradingSecuritiesDebt' xlink:label='loc_TradingSecuritiesDebt' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AvailableForSaleSecuritiesDebtSecuritiesCurrent' xlink:label='loc_AvailableForSaleSecuritiesDebtSecuritiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossCurrentAbstract' xlink:label='loc_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesCurrent' xlink:label='loc_DebtSecuritiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_HeldToMaturitySecuritiesCurrent' xlink:label='loc_HeldToMaturitySecuritiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossNoncurrentAbstract' xlink:label='loc_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_HeldToMaturitySecuritiesNoncurrent' xlink:label='loc_HeldToMaturitySecuritiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityAmortizedCostAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_DebtSecuritiesHeldToMaturityExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesNoncurrentAbstract' xlink:label='loc_DebtSecuritiesNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AvailableForSaleSecuritiesDebtSecuritiesNoncurrent' xlink:label='loc_AvailableForSaleSecuritiesDebtSecuritiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtSecuritiesNoncurrent' xlink:label='loc_DebtSecuritiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtCurrentAbstract' xlink:label='loc_DebtCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermBorrowingsAbstract' xlink:label='loc_ShortTermBorrowingsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtAndCapitalLeaseObligationsCurrentAbstract' xlink:label='loc_LongTermDebtAndCapitalLeaseObligationsCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DebtCurrent' xlink:label='loc_DebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCompensationLiabilityCurrentAbstract' xlink:label='loc_DeferredCompensationLiabilityCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCompensationShareBasedArrangementsLiabilityCurrent' xlink:label='loc_DeferredCompensationShareBasedArrangementsLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCompensationCashBasedArrangementsLiabilityCurrent' xlink:label='loc_DeferredCompensationCashBasedArrangementsLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherDeferredCompensationArrangementsLiabilityCurrent' xlink:label='loc_OtherDeferredCompensationArrangementsLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCompensationLiabilityCurrent' xlink:label='loc_DeferredCompensationLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCostsLeasingNetCurrent' xlink:label='loc_DeferredCostsLeasingNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredFuelCost' xlink:label='loc_DeferredFuelCost' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredGasCost' xlink:label='loc_DeferredGasCost' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredStormAndPropertyReserveDeficiencyCurrent' xlink:label='loc_DeferredStormAndPropertyReserveDeficiencyCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredOfferingCosts' xlink:label='loc_DeferredOfferingCosts' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherDeferredCostsNet' xlink:label='loc_OtherDeferredCostsNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCostsCurrent' xlink:label='loc_DeferredCostsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCostsLeasingNetNoncurrent' xlink:label='loc_DeferredCostsLeasingNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCosts' xlink:label='loc_DeferredCosts' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRevenueAndCreditsCurrentAbstract' xlink:label='loc_DeferredRevenueAndCreditsCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerLiabilityCurrent' xlink:label='loc_ContractWithCustomerLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredIncomeCurrent' xlink:label='loc_DeferredIncomeCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRevenueCurrent' xlink:label='loc_DeferredRevenueCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRevenueAndCreditsNoncurrentAbstract' xlink:label='loc_DeferredRevenueAndCreditsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ContractWithCustomerLiabilityNoncurrent' xlink:label='loc_ContractWithCustomerLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredIncomeNoncurrent' xlink:label='loc_DeferredIncomeNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRevenueNoncurrent' xlink:label='loc_DeferredRevenueNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeAssetsCurrent' xlink:label='loc_DerivativeAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_HedgingAssetsCurrent' xlink:label='loc_HedgingAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommodityContractAssetCurrent' xlink:label='loc_CommodityContractAssetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EnergyMarketingContractsAssetsCurrent' xlink:label='loc_EnergyMarketingContractsAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedges' xlink:label='loc_DerivativeInstrumentsAndHedges' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesLiabilitiesAbstract' xlink:label='loc_DerivativeInstrumentsAndHedgesLiabilitiesAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeLiabilitiesCurrent' xlink:label='loc_DerivativeLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_HedgingLiabilitiesCurrent' xlink:label='loc_HedgingLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EnergyMarketingContractLiabilitiesCurrent' xlink:label='loc_EnergyMarketingContractLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesLiabilities' xlink:label='loc_DerivativeInstrumentsAndHedgesLiabilities' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesLiabilitiesNoncurrentAbstract' xlink:label='loc_DerivativeInstrumentsAndHedgesLiabilitiesNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeLiabilitiesNoncurrent' xlink:label='loc_DerivativeLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_HedgingLiabilitiesNoncurrent' xlink:label='loc_HedgingLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EnergyMarketingContractLiabilitiesNoncurrent' xlink:label='loc_EnergyMarketingContractLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesLiabilitiesNoncurrent' xlink:label='loc_DerivativeInstrumentsAndHedgesLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeAssetsNoncurrent' xlink:label='loc_DerivativeAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_HedgingAssetsNoncurrent' xlink:label='loc_HedgingAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommodityContractAssetNoncurrent' xlink:label='loc_CommodityContractAssetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EnergyMarketingContractsAssetsNoncurrent' xlink:label='loc_EnergyMarketingContractsAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DerivativeInstrumentsAndHedgesNoncurrent' xlink:label='loc_DerivativeInstrumentsAndHedgesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_DirectFinancingLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationCashAndCashEquivalents' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationCashAndCashEquivalents' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationAccountsNotesAndLoansReceivableNet' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationAccountsNotesAndLoansReceivableNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationInventoryCurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationInventoryCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationOtherCurrentAssets' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationOtherCurrentAssets' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationPrepaidAndOtherAssetsCurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationPrepaidAndOtherAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationGoodwillCurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationGoodwillCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationPropertyPlantAndEquipmentCurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationPropertyPlantAndEquipmentCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationIntangibleAssetsCurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationIntangibleAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetsOfDisposalGroupIncludingDiscontinuedOperationCurrent' xlink:label='loc_AssetsOfDisposalGroupIncludingDiscontinuedOperationCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationOtherNoncurrentAssets' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationOtherNoncurrentAssets' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationInventoryNoncurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationInventoryNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationDeferredTaxAssets' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationDeferredTaxAssets' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationIntangibleAssetsNoncurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationIntangibleAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationGoodwillNoncurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationGoodwillNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationPropertyPlantAndEquipmentNoncurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationPropertyPlantAndEquipmentNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DisposalGroupIncludingDiscontinuedOperationAssetsNoncurrent' xlink:label='loc_DisposalGroupIncludingDiscontinuedOperationAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:label='loc_FinancingReceivableExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:label='loc_FinancingReceivableAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_FinancingReceivableExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_FinancingReceivableAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_FinancingReceivableExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FiniteLivedIntangibleAssetsNetAbstract' xlink:label='loc_FiniteLivedIntangibleAssetsNetAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FiniteLivedIntangibleAssetsGross' xlink:label='loc_FiniteLivedIntangibleAssetsGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FiniteLivedIntangibleAssetsAccumulatedAmortization' xlink:label='loc_FiniteLivedIntangibleAssetsAccumulatedAmortization' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FiniteLivedIntangibleAssetsNet' xlink:label='loc_FiniteLivedIntangibleAssetsNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_GeneralPartnersCapitalAccountAbstract' xlink:label='loc_GeneralPartnersCapitalAccountAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_GeneralPartnersCapitalAccountUnitsAuthorized' xlink:label='loc_GeneralPartnersCapitalAccountUnitsAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_GeneralPartnersCapitalAccountUnitsIssued' xlink:label='loc_GeneralPartnersCapitalAccountUnitsIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_GeneralPartnersCapitalAccountUnitsOutstanding' xlink:label='loc_GeneralPartnersCapitalAccountUnitsOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_IndefiniteLivedIntangibleAssetsExcludingGoodwill' xlink:label='loc_IndefiniteLivedIntangibleAssetsExcludingGoodwill' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_IntangibleAssetsNetExcludingGoodwill' xlink:label='loc_IntangibleAssetsNetExcludingGoodwill' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryFinishedGoods' xlink:label='loc_InventoryFinishedGoods' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryForLongTermContractsOrPrograms' xlink:label='loc_InventoryForLongTermContractsOrPrograms' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryWorkInProcess' xlink:label='loc_InventoryWorkInProcess' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryRawMaterialsAndSupplies' xlink:label='loc_InventoryRawMaterialsAndSupplies' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryValuationReserves' xlink:label='loc_InventoryValuationReserves' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryLIFOReserve' xlink:label='loc_InventoryLIFOReserve' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryNet' xlink:label='loc_InventoryNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryGasInStorageUndergroundNoncurrent' xlink:label='loc_InventoryGasInStorageUndergroundNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryDrillingNoncurrent' xlink:label='loc_InventoryDrillingNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherInventoryNoncurrent' xlink:label='loc_OtherInventoryNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InventoryNoncurrent' xlink:label='loc_InventoryNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InvestmentsInAffiliatesSubsidiariesAssociatesAndJointVenturesAbstract' xlink:label='loc_InvestmentsInAffiliatesSubsidiariesAssociatesAndJointVenturesAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EquityMethodInvestments' xlink:label='loc_EquityMethodInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AdvancesToAffiliate' xlink:label='loc_AdvancesToAffiliate' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InvestmentsInAffiliatesSubsidiariesAssociatesAndJointVentures' xlink:label='loc_InvestmentsInAffiliatesSubsidiariesAssociatesAndJointVentures' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LandAndLandImprovementsAbstract' xlink:label='loc_LandAndLandImprovementsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_Land' xlink:label='loc_Land' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LandImprovements' xlink:label='loc_LandImprovements' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LandAndLandImprovements' xlink:label='loc_LandAndLandImprovements' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesAbstract' xlink:label='loc_LiabilitiesAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesCurrentAbstract' xlink:label='loc_LiabilitiesCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesNoncurrentAbstract' xlink:label='loc_LiabilitiesNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_Liabilities' xlink:label='loc_Liabilities' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesAndStockholdersEquityAbstract' xlink:label='loc_LiabilitiesAndStockholdersEquityAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommitmentsAndContingencies' xlink:label='loc_CommitmentsAndContingencies' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityAbstract' xlink:label='loc_TemporaryEquityAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StockholdersEquityIncludingPortionAttributableToNoncontrollingInterestAbstract' xlink:label='loc_StockholdersEquityIncludingPortionAttributableToNoncontrollingInterestAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StockholdersEquityNumberOfSharesParValueAndOtherDisclosuresAbstract' xlink:label='loc_StockholdersEquityNumberOfSharesParValueAndOtherDisclosuresAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapitalIncludingPortionAttributableToNoncontrollingInterestAbstract' xlink:label='loc_PartnersCapitalIncludingPortionAttributableToNoncontrollingInterestAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapitalNumberOfUnitsParValueAndOtherDisclosuresAbstract' xlink:label='loc_PartnersCapitalNumberOfUnitsParValueAndOtherDisclosuresAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLLCMembersEquityIncludingPortionAttributableToNoncontrollingInterestAbstract' xlink:label='loc_LimitedLiabilityCompanyLLCMembersEquityIncludingPortionAttributableToNoncontrollingInterestAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesAndStockholdersEquity' xlink:label='loc_LiabilitiesAndStockholdersEquity' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRentCreditCurrent' xlink:label='loc_DeferredRentCreditCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilityForUncertainTaxPositionsCurrent' xlink:label='loc_LiabilityForUncertainTaxPositionsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PostemploymentBenefitsLiabilityCurrent' xlink:label='loc_PostemploymentBenefitsLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SecuritiesLoaned' xlink:label='loc_SecuritiesLoaned' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RegulatoryLiabilityCurrent' xlink:label='loc_RegulatoryLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ProvisionForLossOnContracts' xlink:label='loc_ProvisionForLossOnContracts' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LitigationReserveCurrent' xlink:label='loc_LitigationReserveCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedEnvironmentalLossContingenciesCurrent' xlink:label='loc_AccruedEnvironmentalLossContingenciesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetRetirementObligationCurrent' xlink:label='loc_AssetRetirementObligationCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedCappingClosurePostClosureAndEnvironmentalCosts' xlink:label='loc_AccruedCappingClosurePostClosureAndEnvironmentalCosts' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedReclamationCostsCurrent' xlink:label='loc_AccruedReclamationCostsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ProgramRightsObligationsCurrent' xlink:label='loc_ProgramRightsObligationsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DueToRelatedPartiesCurrent' xlink:label='loc_DueToRelatedPartiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesOfDisposalGroupIncludingDiscontinuedOperationCurrent' xlink:label='loc_LiabilitiesOfDisposalGroupIncludingDiscontinuedOperationCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredGasPurchasesCurrent' xlink:label='loc_DeferredGasPurchasesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesOfBusinessTransferredUnderContractualArrangementCurrent' xlink:label='loc_LiabilitiesOfBusinessTransferredUnderContractualArrangementCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLiabilitiesCurrent' xlink:label='loc_OtherLiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CustomerRefundLiabilityCurrent' xlink:label='loc_CustomerRefundLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SelfInsuranceReserveCurrent' xlink:label='loc_SelfInsuranceReserveCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetAcquisitionContingentConsiderationLiabilityCurrent' xlink:label='loc_AssetAcquisitionContingentConsiderationLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BusinessCombinationContingentConsiderationLiabilityCurrent' xlink:label='loc_BusinessCombinationContingentConsiderationLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EntertainmentLicenseAgreementForProgramMaterialLiabilityCurrent' xlink:label='loc_EntertainmentLicenseAgreementForProgramMaterialLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesCurrent' xlink:label='loc_LiabilitiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtAndCapitalLeaseObligationsAbstract' xlink:label='loc_LongTermDebtAndCapitalLeaseObligationsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesOtherThanLongTermDebtNoncurrentAbstract' xlink:label='loc_LiabilitiesOtherThanLongTermDebtNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesNoncurrent' xlink:label='loc_LiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsPayableAndAccruedLiabilitiesNoncurrent' xlink:label='loc_AccountsPayableAndAccruedLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCompensationLiabilityClassifiedNoncurrent' xlink:label='loc_DeferredCompensationLiabilityClassifiedNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PensionAndOtherPostretirementDefinedBenefitPlansLiabilitiesNoncurrentAbstract' xlink:label='loc_PensionAndOtherPostretirementDefinedBenefitPlansLiabilitiesNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedDeferredInvestmentTaxCredit' xlink:label='loc_AccumulatedDeferredInvestmentTaxCredit' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredRentCreditNoncurrent' xlink:label='loc_DeferredRentCreditNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredIncomeTaxLiabilitiesNet' xlink:label='loc_DeferredIncomeTaxLiabilitiesNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilityForUncertainTaxPositionsNoncurrent' xlink:label='loc_LiabilityForUncertainTaxPositionsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PostemploymentBenefitsLiabilityNoncurrent' xlink:label='loc_PostemploymentBenefitsLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedEnvironmentalLossContingenciesNoncurrent' xlink:label='loc_AccruedEnvironmentalLossContingenciesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CustomerRefundLiabilityNoncurrent' xlink:label='loc_CustomerRefundLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OffMarketLeaseUnfavorable' xlink:label='loc_OffMarketLeaseUnfavorable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LeaseDepositLiability' xlink:label='loc_LeaseDepositLiability' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SharesSubjectToMandatoryRedemptionSettlementTermsAmountNoncurrent' xlink:label='loc_SharesSubjectToMandatoryRedemptionSettlementTermsAmountNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LitigationReserveNoncurrent' xlink:label='loc_LitigationReserveNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RegulatoryLiabilityNoncurrent' xlink:label='loc_RegulatoryLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestructuringReserveNoncurrent' xlink:label='loc_RestructuringReserveNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DueToRelatedPartiesNoncurrent' xlink:label='loc_DueToRelatedPartiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesOfDisposalGroupIncludingDiscontinuedOperationNoncurrent' xlink:label='loc_LiabilitiesOfDisposalGroupIncludingDiscontinuedOperationNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ProgramRightsObligationsNoncurrent' xlink:label='loc_ProgramRightsObligationsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesOfBusinessTransferredUnderContractualArrangementNoncurrent' xlink:label='loc_LiabilitiesOfBusinessTransferredUnderContractualArrangementNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SelfInsuranceReserveNoncurrent' xlink:label='loc_SelfInsuranceReserveNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLiabilitiesNoncurrent' xlink:label='loc_OtherLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetAcquisitionContingentConsiderationLiabilityNoncurrent' xlink:label='loc_AssetAcquisitionContingentConsiderationLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BusinessCombinationContingentConsiderationLiabilityNoncurrent' xlink:label='loc_BusinessCombinationContingentConsiderationLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EntertainmentLicenseAgreementForProgramMaterialLiabilityNoncurrent' xlink:label='loc_EntertainmentLicenseAgreementForProgramMaterialLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OperatingLeaseLiabilityNoncurrent' xlink:label='loc_OperatingLeaseLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_QualifiedAffordableHousingProjectInvestmentsCommitment' xlink:label='loc_QualifiedAffordableHousingProjectInvestmentsCommitment' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LiabilitiesOtherThanLongtermDebtNoncurrent' xlink:label='loc_LiabilitiesOtherThanLongtermDebtNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DefinedBenefitPensionPlanLiabilitiesNoncurrent' xlink:label='loc_DefinedBenefitPensionPlanLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherPostretirementDefinedBenefitPlanLiabilitiesNoncurrent' xlink:label='loc_OtherPostretirementDefinedBenefitPlanLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PensionAndOtherPostretirementDefinedBenefitPlansLiabilitiesNoncurrent' xlink:label='loc_PensionAndOtherPostretirementDefinedBenefitPlansLiabilitiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLLCMembersEquityAbstract' xlink:label='loc_LimitedLiabilityCompanyLLCMembersEquityAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonUnitIssued' xlink:label='loc_CommonUnitIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonUnitAuthorized' xlink:label='loc_CommonUnitAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonUnitOutstanding' xlink:label='loc_CommonUnitOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonUnitIssuanceValue' xlink:label='loc_CommonUnitIssuanceValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLLCPreferredUnitIssued' xlink:label='loc_LimitedLiabilityCompanyLLCPreferredUnitIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLLCPreferredUnitAuthorized' xlink:label='loc_LimitedLiabilityCompanyLLCPreferredUnitAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLLCPreferredUnitOutstanding' xlink:label='loc_LimitedLiabilityCompanyLLCPreferredUnitOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLLCPreferredUnitIssuanceValue' xlink:label='loc_LimitedLiabilityCompanyLLCPreferredUnitIssuanceValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MembersEquityAbstract' xlink:label='loc_MembersEquityAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MembersEquityAttributableToNoncontrollingInterest' xlink:label='loc_MembersEquityAttributableToNoncontrollingInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedLiabilityCompanyLlcMembersEquityIncludingPortionAttributableToNoncontrollingInterest' xlink:label='loc_LimitedLiabilityCompanyLlcMembersEquityIncludingPortionAttributableToNoncontrollingInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedPartnersCapitalAccountAbstract' xlink:label='loc_LimitedPartnersCapitalAccountAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedPartnersCapitalAccountUnitsAuthorized' xlink:label='loc_LimitedPartnersCapitalAccountUnitsAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedPartnersCapitalAccountUnitsIssued' xlink:label='loc_LimitedPartnersCapitalAccountUnitsIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedPartnersCapitalAccountUnitsOutstanding' xlink:label='loc_LimitedPartnersCapitalAccountUnitsOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LoansPayableCurrentAbstract' xlink:label='loc_LoansPayableCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LoansPayableToBankCurrent' xlink:label='loc_LoansPayableToBankCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLoansPayableCurrent' xlink:label='loc_OtherLoansPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LoansPayableCurrent' xlink:label='loc_LoansPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermLoansPayableAbstract' xlink:label='loc_LongTermLoansPayableAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermLoansFromBank' xlink:label='loc_LongTermLoansFromBank' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLoansPayableLongTerm' xlink:label='loc_OtherLoansPayableLongTerm' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermLoansPayable' xlink:label='loc_LongTermLoansPayable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtCurrentAbstract' xlink:label='loc_LongTermDebtCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SecuredDebtCurrent' xlink:label='loc_SecuredDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertibleDebtCurrent' xlink:label='loc_ConvertibleDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_UnsecuredDebtCurrent' xlink:label='loc_UnsecuredDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SubordinatedDebtCurrent' xlink:label='loc_SubordinatedDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertibleSubordinatedDebtCurrent' xlink:label='loc_ConvertibleSubordinatedDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermCommercialPaperCurrent' xlink:label='loc_LongTermCommercialPaperCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermConstructionLoanCurrent' xlink:label='loc_LongTermConstructionLoanCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongtermTransitionBondCurrent' xlink:label='loc_LongtermTransitionBondCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongtermPollutionControlBondCurrent' xlink:label='loc_LongtermPollutionControlBondCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_JuniorSubordinatedDebentureOwedToUnconsolidatedSubsidiaryTrustCurrent' xlink:label='loc_JuniorSubordinatedDebentureOwedToUnconsolidatedSubsidiaryTrustCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLongTermDebtCurrent' xlink:label='loc_OtherLongTermDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LinesOfCreditCurrent' xlink:label='loc_LinesOfCreditCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansPayableCurrentAbstract' xlink:label='loc_NotesAndLoansPayableCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SpecialAssessmentBondCurrent' xlink:label='loc_SpecialAssessmentBondCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FederalHomeLoanBankAdvancesCurrent' xlink:label='loc_FederalHomeLoanBankAdvancesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtCurrent' xlink:label='loc_LongTermDebtCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtCurrentRecourseStatusExtensibleEnumeration' xlink:label='loc_LongTermDebtCurrentRecourseStatusExtensibleEnumeration' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtNoncurrentAbstract' xlink:label='loc_LongTermDebtNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermLineOfCredit' xlink:label='loc_LongTermLineOfCredit' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommercialPaperNoncurrent' xlink:label='loc_CommercialPaperNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConstructionLoanNoncurrent' xlink:label='loc_ConstructionLoanNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SecuredLongTermDebt' xlink:label='loc_SecuredLongTermDebt' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SubordinatedLongTermDebt' xlink:label='loc_SubordinatedLongTermDebt' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_UnsecuredLongTermDebt' xlink:label='loc_UnsecuredLongTermDebt' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertibleDebtNoncurrent' xlink:label='loc_ConvertibleDebtNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertibleSubordinatedDebtNoncurrent' xlink:label='loc_ConvertibleSubordinatedDebtNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermTransitionBond' xlink:label='loc_LongTermTransitionBond' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermPollutionControlBond' xlink:label='loc_LongTermPollutionControlBond' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_JuniorSubordinatedDebentureOwedToUnconsolidatedSubsidiaryTrustNoncurrent' xlink:label='loc_JuniorSubordinatedDebentureOwedToUnconsolidatedSubsidiaryTrustNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermNotesAndLoansAbstract' xlink:label='loc_LongTermNotesAndLoansAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SpecialAssessmentBondNoncurrent' xlink:label='loc_SpecialAssessmentBondNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongtermFederalHomeLoanBankAdvancesNoncurrent' xlink:label='loc_LongtermFederalHomeLoanBankAdvancesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLongTermDebtNoncurrent' xlink:label='loc_OtherLongTermDebtNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtNoncurrent' xlink:label='loc_LongTermDebtNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtNoncurrentRecourseStatusExtensibleEnumeration' xlink:label='loc_LongTermDebtNoncurrentRecourseStatusExtensibleEnumeration' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinanceLeaseLiabilityNoncurrent' xlink:label='loc_FinanceLeaseLiabilityNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtAndCapitalLeaseObligations' xlink:label='loc_LongTermDebtAndCapitalLeaseObligations' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FinanceLeaseLiabilityCurrent' xlink:label='loc_FinanceLeaseLiabilityCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermDebtAndCapitalLeaseObligationsCurrent' xlink:label='loc_LongTermDebtAndCapitalLeaseObligationsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermInvestmentsAbstract' xlink:label='loc_LongTermInvestmentsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EquitySecuritiesFVNINoncurrent' xlink:label='loc_EquitySecuritiesFVNINoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MarketableSecuritiesNoncurrent' xlink:label='loc_MarketableSecuritiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InvestmentInPhysicalCommodities' xlink:label='loc_InvestmentInPhysicalCommodities' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLongTermInvestments' xlink:label='loc_OtherLongTermInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AuctionRateSecuritiesNoncurrent' xlink:label='loc_AuctionRateSecuritiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermInvestments' xlink:label='loc_LongTermInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DueFromRelatedPartiesNoncurrent' xlink:label='loc_DueFromRelatedPartiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermInvestmentsAndReceivablesNet' xlink:label='loc_LongTermInvestmentsAndReceivablesNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MembersCapital' xlink:label='loc_MembersCapital' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RetainedEarningsAccumulatedDeficit' xlink:label='loc_RetainedEarningsAccumulatedDeficit' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesReceivableByOwnerToLimitedLiabilityCompanyLLC' xlink:label='loc_NotesReceivableByOwnerToLimitedLiabilityCompanyLLC' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MembersEquity' xlink:label='loc_MembersEquity' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:label='loc_NetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:label='loc_NetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_NetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_NetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_NetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesPayableCurrentAbstract' xlink:label='loc_NotesPayableCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MediumtermNotesCurrent' xlink:label='loc_MediumtermNotesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertibleNotesPayableCurrent' xlink:label='loc_ConvertibleNotesPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesPayableToBankCurrent' xlink:label='loc_NotesPayableToBankCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeniorNotesCurrent' xlink:label='loc_SeniorNotesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_JuniorSubordinatedNotesCurrent' xlink:label='loc_JuniorSubordinatedNotesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherNotesPayableCurrent' xlink:label='loc_OtherNotesPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesPayableCurrent' xlink:label='loc_NotesPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermNotesPayableAbstract' xlink:label='loc_LongTermNotesPayableAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MediumtermNotesNoncurrent' xlink:label='loc_MediumtermNotesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_JuniorSubordinatedLongTermNotes' xlink:label='loc_JuniorSubordinatedLongTermNotes' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SeniorLongTermNotes' xlink:label='loc_SeniorLongTermNotes' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConvertibleLongTermNotesPayable' xlink:label='loc_ConvertibleLongTermNotesPayable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesPayableToBankNoncurrent' xlink:label='loc_NotesPayableToBankNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherLongTermNotesPayable' xlink:label='loc_OtherLongTermNotesPayable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermNotesPayable' xlink:label='loc_LongTermNotesPayable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansPayableCurrent' xlink:label='loc_NotesAndLoansPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LongTermNotesAndLoans' xlink:label='loc_LongTermNotesAndLoans' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansReceivableGrossCurrent' xlink:label='loc_NotesAndLoansReceivableGrossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AllowanceForNotesAndLoansReceivableCurrent' xlink:label='loc_AllowanceForNotesAndLoansReceivableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansReceivableNetCurrent' xlink:label='loc_NotesAndLoansReceivableNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansReceivableGrossNoncurrent' xlink:label='loc_NotesAndLoansReceivableGrossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AllowanceForNotesAndLoansReceivableNoncurrent' xlink:label='loc_AllowanceForNotesAndLoansReceivableNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NotesAndLoansReceivableNetNoncurrent' xlink:label='loc_NotesAndLoansReceivableNetNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherUnitsOtherOwnershipInterestsCapitalAccountAbstract' xlink:label='loc_OtherUnitsOtherOwnershipInterestsCapitalAccountAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherOwnershipInterestsUnitsAuthorized' xlink:label='loc_OtherOwnershipInterestsUnitsAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherOwnershipInterestsUnitsIssued' xlink:label='loc_OtherOwnershipInterestsUnitsIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherOwnershipInterestsUnitsOutstanding' xlink:label='loc_OtherOwnershipInterestsUnitsOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapitalAbstract' xlink:label='loc_PartnersCapitalAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_GeneralPartnersCapitalAccount' xlink:label='loc_GeneralPartnersCapitalAccount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LimitedPartnersCapitalAccount' xlink:label='loc_LimitedPartnersCapitalAccount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredUnitsPreferredPartnersCapitalAccounts' xlink:label='loc_PreferredUnitsPreferredPartnersCapitalAccounts' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherPartnersCapital' xlink:label='loc_OtherPartnersCapital' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapitalAllocatedForIncomeTaxAndOtherWithdrawals' xlink:label='loc_PartnersCapitalAllocatedForIncomeTaxAndOtherWithdrawals' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherOwnershipInterestsCapitalAccount' xlink:label='loc_OtherOwnershipInterestsCapitalAccount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OfferingCostsPartnershipInterests' xlink:label='loc_OfferingCostsPartnershipInterests' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapital' xlink:label='loc_PartnersCapital' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapitalAttributableToNoncontrollingInterest' xlink:label='loc_PartnersCapitalAttributableToNoncontrollingInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PartnersCapitalIncludingPortionAttributableToNoncontrollingInterest' xlink:label='loc_PartnersCapitalIncludingPortionAttributableToNoncontrollingInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredUnitsPreferredPartnersCapitalAccountAbstract' xlink:label='loc_PreferredUnitsPreferredPartnersCapitalAccountAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockNumberOfSharesParValueAndOtherDisclosuresAbstract' xlink:label='loc_PreferredStockNumberOfSharesParValueAndOtherDisclosuresAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockParOrStatedValuePerShare' xlink:label='loc_PreferredStockParOrStatedValuePerShare' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockNoParValue' xlink:label='loc_PreferredStockNoParValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockSharesSubscribedButUnissuedValue' xlink:label='loc_PreferredStockSharesSubscribedButUnissuedValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockShareSubscriptions' xlink:label='loc_PreferredStockShareSubscriptions' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockSharesAuthorized' xlink:label='loc_PreferredStockSharesAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockSharesAuthorizedUnlimited' xlink:label='loc_PreferredStockSharesAuthorizedUnlimited' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockSharesIssued' xlink:label='loc_PreferredStockSharesIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockSharesOutstanding' xlink:label='loc_PreferredStockSharesOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockValueOutstanding' xlink:label='loc_PreferredStockValueOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockRedemptionAmount' xlink:label='loc_PreferredStockRedemptionAmount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockRedemptionAmountFutureRedeemableSecurities' xlink:label='loc_PreferredStockRedemptionAmountFutureRedeemableSecurities' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockLiquidationPreferenceValue' xlink:label='loc_PreferredStockLiquidationPreferenceValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockLiquidationPreference' xlink:label='loc_PreferredStockLiquidationPreference' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockConversionBasis' xlink:label='loc_PreferredStockConversionBasis' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredUnitsAuthorized' xlink:label='loc_PreferredUnitsAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredUnitsIssued' xlink:label='loc_PreferredUnitsIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredUnitsOutstanding' xlink:label='loc_PreferredUnitsOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidInsurance' xlink:label='loc_PrepaidInsurance' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidRent' xlink:label='loc_PrepaidRent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidAdvertising' xlink:label='loc_PrepaidAdvertising' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidRoyalties' xlink:label='loc_PrepaidRoyalties' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_Supplies' xlink:label='loc_Supplies' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidInterest' xlink:label='loc_PrepaidInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidTaxes' xlink:label='loc_PrepaidTaxes' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherPrepaidExpenseCurrent' xlink:label='loc_OtherPrepaidExpenseCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidExpenseCurrent' xlink:label='loc_PrepaidExpenseCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidExpenseOtherNoncurrent' xlink:label='loc_PrepaidExpenseOtherNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidMineralRoyaltiesNoncurrent' xlink:label='loc_PrepaidMineralRoyaltiesNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PrepaidExpenseNoncurrent' xlink:label='loc_PrepaidExpenseNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PropertyPlantAndEquipmentGrossAbstract' xlink:label='loc_PropertyPlantAndEquipmentGrossAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BuildingsAndImprovementsGross' xlink:label='loc_BuildingsAndImprovementsGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MachineryAndEquipmentGross' xlink:label='loc_MachineryAndEquipmentGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_FurnitureAndFixturesGross' xlink:label='loc_FurnitureAndFixturesGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CapitalizedComputerSoftwareGross' xlink:label='loc_CapitalizedComputerSoftwareGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConstructionInProgressGross' xlink:label='loc_ConstructionInProgressGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_LeaseholdImprovementsGross' xlink:label='loc_LeaseholdImprovementsGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TimberAndTimberlands' xlink:label='loc_TimberAndTimberlands' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PropertyPlantAndEquipmentOther' xlink:label='loc_PropertyPlantAndEquipmentOther' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PropertyPlantAndEquipmentGross' xlink:label='loc_PropertyPlantAndEquipmentGross' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccumulatedDepreciationDepletionAndAmortizationPropertyPlantAndEquipment' xlink:label='loc_AccumulatedDepreciationDepletionAndAmortizationPropertyPlantAndEquipment' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PropertyPlantAndEquipmentNet' xlink:label='loc_PropertyPlantAndEquipmentNet' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_NontradeReceivablesCurrent' xlink:label='loc_NontradeReceivablesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_UnbilledReceivablesCurrent' xlink:label='loc_UnbilledReceivablesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DueFromRelatedPartiesCurrent' xlink:label='loc_DueFromRelatedPartiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ReceivablesLongTermContractsOrPrograms' xlink:label='loc_ReceivablesLongTermContractsOrPrograms' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsReceivableFromSecuritization' xlink:label='loc_AccountsReceivableFromSecuritization' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccountsAndOtherReceivablesNetCurrent' xlink:label='loc_AccountsAndOtherReceivablesNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ReceivablesNetCurrent' xlink:label='loc_ReceivablesNetCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemableNoncontrollingInterestEquityCarryingAmountAbstract' xlink:label='loc_RedeemableNoncontrollingInterestEquityCarryingAmountAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemableNoncontrollingInterestEquityCommonCarryingAmount' xlink:label='loc_RedeemableNoncontrollingInterestEquityCommonCarryingAmount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemableNoncontrollingInterestEquityPreferredCarryingAmount' xlink:label='loc_RedeemableNoncontrollingInterestEquityPreferredCarryingAmount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemableNoncontrollingInterestEquityOtherCarryingAmount' xlink:label='loc_RedeemableNoncontrollingInterestEquityOtherCarryingAmount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RedeemableNoncontrollingInterestEquityCarryingAmount' xlink:label='loc_RedeemableNoncontrollingInterestEquityCarryingAmount' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RegulatoryAssetsNoncurrent' xlink:label='loc_RegulatoryAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SecuritizedRegulatoryTransitionAssetsNoncurrent' xlink:label='loc_SecuritizedRegulatoryTransitionAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DemandSideManagementProgramCostsNoncurrent' xlink:label='loc_DemandSideManagementProgramCostsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_UnamortizedLossReacquiredDebtNoncurrent' xlink:label='loc_UnamortizedLossReacquiredDebtNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DecommissioningFundInvestments' xlink:label='loc_DecommissioningFundInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AssetRecoveryDamagedPropertyCostsNoncurrent' xlink:label='loc_AssetRecoveryDamagedPropertyCostsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredStormAndPropertyReserveDeficiencyNoncurrent' xlink:label='loc_DeferredStormAndPropertyReserveDeficiencyNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_InvestmentsInPowerAndDistributionProjects' xlink:label='loc_InvestmentsInPowerAndDistributionProjects' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PhaseInPlanAmountOfCostsDeferredForRateMakingPurposes' xlink:label='loc_PhaseInPlanAmountOfCostsDeferredForRateMakingPurposes' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_UnamortizedDebtIssuanceExpense' xlink:label='loc_UnamortizedDebtIssuanceExpense' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RegulatedEntityOtherAssetsNoncurrent' xlink:label='loc_RegulatedEntityOtherAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='https://xbrl.fasb.org/srt/2022/elts/srt-2022.xsd#srt_RestatementAxis' xlink:label='loc_RestatementAxis' xlink:type='locator' />
    <link:loc xlink:href='https://xbrl.fasb.org/srt/2022/elts/srt-2022.xsd#srt_RestatementDomain' xlink:label='loc_RestatementDomain' xlink:type='locator' />
    <link:loc xlink:href='https://xbrl.fasb.org/srt/2022/elts/srt-2022.xsd#srt_ScenarioPreviouslyReportedMember' xlink:label='loc_ScenarioPreviouslyReportedMember' xlink:type='locator' />
    <link:loc xlink:href='https://xbrl.fasb.org/srt/2022/elts/srt-2022.xsd#srt_RestatementAdjustmentMember' xlink:label='loc_RestatementAdjustmentMember' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndCashEquivalentsCurrentAbstract' xlink:label='loc_RestrictedCashAndCashEquivalentsCurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashCurrent' xlink:label='loc_RestrictedCashCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashEquivalentsCurrent' xlink:label='loc_RestrictedCashEquivalentsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndCashEquivalentsAtCarryingValue' xlink:label='loc_RestrictedCashAndCashEquivalentsAtCarryingValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndCashEquivalentsNoncurrentAbstract' xlink:label='loc_RestrictedCashAndCashEquivalentsNoncurrentAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashNoncurrent' xlink:label='loc_RestrictedCashNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashEquivalentsNoncurrent' xlink:label='loc_RestrictedCashEquivalentsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndCashEquivalentsNoncurrent' xlink:label='loc_RestrictedCashAndCashEquivalentsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedInvestmentsCurrent' xlink:label='loc_RestrictedInvestmentsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherRestrictedAssetsCurrent' xlink:label='loc_OtherRestrictedAssetsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndInvestmentsCurrent' xlink:label='loc_RestrictedCashAndInvestmentsCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedInvestmentsNoncurrent' xlink:label='loc_RestrictedInvestmentsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherRestrictedAssetsNoncurrent' xlink:label='loc_OtherRestrictedAssetsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RestrictedCashAndInvestmentsNoncurrent' xlink:label='loc_RestrictedCashAndInvestmentsNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RetainedEarningsAccumulatedDeficitAbstract' xlink:label='loc_RetainedEarningsAccumulatedDeficitAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RetainedEarningsAppropriated' xlink:label='loc_RetainedEarningsAppropriated' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RetainedEarningsUnappropriated' xlink:label='loc_RetainedEarningsUnappropriated' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestBeforeAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseAllowanceForCreditLossExcludingAccruedInterestNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:label='loc_SalesTypeLeaseNetInvestmentInLeaseExcludingAccruedInterestAfterAllowanceForCreditLossNoncurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BankOverdrafts' xlink:label='loc_BankOverdrafts' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommercialPaper' xlink:label='loc_CommercialPaper' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_BridgeLoan' xlink:label='loc_BridgeLoan' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ConstructionLoan' xlink:label='loc_ConstructionLoan' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermBankLoansAndNotesPayable' xlink:label='loc_ShortTermBankLoansAndNotesPayable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermNonBankLoansAndNotesPayable' xlink:label='loc_ShortTermNonBankLoansAndNotesPayable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherShortTermBorrowings' xlink:label='loc_OtherShortTermBorrowings' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermBorrowings' xlink:label='loc_ShortTermBorrowings' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermDebtRecourseStatusExtensibleEnumeration' xlink:label='loc_ShortTermDebtRecourseStatusExtensibleEnumeration' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_EquitySecuritiesFvNi' xlink:label='loc_EquitySecuritiesFvNi' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MarketableSecuritiesCurrent' xlink:label='loc_MarketableSecuritiesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherShortTermInvestments' xlink:label='loc_OtherShortTermInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ShortTermInvestments' xlink:label='loc_ShortTermInvestments' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StatementLineItems' xlink:label='loc_StatementLineItems' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StatementTable' xlink:label='loc_StatementTable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StatementOfFinancialPositionAbstract' xlink:label='loc_StatementOfFinancialPositionAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StockholdersEquityAbstract' xlink:label='loc_StockholdersEquityAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockValue' xlink:label='loc_PreferredStockValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_PreferredStockSharesSubscribedButUnissuedSubscriptionsReceivable' xlink:label='loc_PreferredStockSharesSubscribedButUnissuedSubscriptionsReceivable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockValue' xlink:label='loc_CommonStockValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockHeldBySubsidiary' xlink:label='loc_CommonStockHeldBySubsidiary' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_CommonStockShareSubscribedButUnissuedSubscriptionsReceivable' xlink:label='loc_CommonStockShareSubscribedButUnissuedSubscriptionsReceivable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_OtherAdditionalCapital' xlink:label='loc_OtherAdditionalCapital' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockValueAbstract' xlink:label='loc_TreasuryStockValueAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockDeferredEmployeeStockOwnershipPlan' xlink:label='loc_TreasuryStockDeferredEmployeeStockOwnershipPlan' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_DeferredCompensationEquity' xlink:label='loc_DeferredCompensationEquity' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_UnearnedESOPShares' xlink:label='loc_UnearnedESOPShares' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ReceivableFromOfficersAndDirectorsForIssuanceOfCapitalStock' xlink:label='loc_ReceivableFromOfficersAndDirectorsForIssuanceOfCapitalStock' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_ReceivableFromShareholdersOrAffiliatesForIssuanceOfCapitalStock' xlink:label='loc_ReceivableFromShareholdersOrAffiliatesForIssuanceOfCapitalStock' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StockholdersEquityNoteSubscriptionsReceivable' xlink:label='loc_StockholdersEquityNoteSubscriptionsReceivable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_WarrantsAndRightsOutstanding' xlink:label='loc_WarrantsAndRightsOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StockholdersEquity' xlink:label='loc_StockholdersEquity' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_MinorityInterest' xlink:label='loc_MinorityInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_StockholdersEquityIncludingPortionAttributableToNoncontrollingInterest' xlink:label='loc_StockholdersEquityIncludingPortionAttributableToNoncontrollingInterest' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockNumberOfSharesAndRestrictionDisclosuresAbstract' xlink:label='loc_TreasuryStockNumberOfSharesAndRestrictionDisclosuresAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_RetainedEarningsDeficitEliminated' xlink:label='loc_RetainedEarningsDeficitEliminated' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccruedIncomeTaxesCurrent' xlink:label='loc_AccruedIncomeTaxesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_SalesAndExciseTaxPayableCurrent' xlink:label='loc_SalesAndExciseTaxPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_AccrualForTaxesOtherThanIncomeTaxesCurrent' xlink:label='loc_AccrualForTaxesOtherThanIncomeTaxesCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TaxesPayableCurrent' xlink:label='loc_TaxesPayableCurrent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityCarryingAmountIncludingPortionAttributableToNoncontrollingInterestsAbstract' xlink:label='loc_TemporaryEquityCarryingAmountIncludingPortionAttributableToNoncontrollingInterestsAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityNumberOfSharesRedemptionValueAndOtherDisclosuresAbstract' xlink:label='loc_TemporaryEquityNumberOfSharesRedemptionValueAndOtherDisclosuresAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityCarryingAmountAttributableToParent' xlink:label='loc_TemporaryEquityCarryingAmountAttributableToParent' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityCarryingAmountIncludingPortionAttributableToNoncontrollingInterests' xlink:label='loc_TemporaryEquityCarryingAmountIncludingPortionAttributableToNoncontrollingInterests' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityParOrStatedValuePerShare' xlink:label='loc_TemporaryEquityParOrStatedValuePerShare' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityValueExcludingAdditionalPaidInCapital' xlink:label='loc_TemporaryEquityValueExcludingAdditionalPaidInCapital' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquitySharesAuthorized' xlink:label='loc_TemporaryEquitySharesAuthorized' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquitySharesIssued' xlink:label='loc_TemporaryEquitySharesIssued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquitySharesOutstanding' xlink:label='loc_TemporaryEquitySharesOutstanding' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityRedemptionPricePerShare' xlink:label='loc_TemporaryEquityRedemptionPricePerShare' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityAggregateAmountOfRedemptionRequirement' xlink:label='loc_TemporaryEquityAggregateAmountOfRedemptionRequirement' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityLiquidationPreferencePerShare' xlink:label='loc_TemporaryEquityLiquidationPreferencePerShare' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityLiquidationPreference' xlink:label='loc_TemporaryEquityLiquidationPreference' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquityShareSubscriptions' xlink:label='loc_TemporaryEquityShareSubscriptions' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquitySharesSubscribedButUnissued' xlink:label='loc_TemporaryEquitySharesSubscribedButUnissued' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TemporaryEquitySharesSubscribedButUnissuedSubscriptionsReceivable' xlink:label='loc_TemporaryEquitySharesSubscribedButUnissuedSubscriptionsReceivable' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockSharesAbstract' xlink:label='loc_TreasuryStockSharesAbstract' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockRestrictions' xlink:label='loc_TreasuryStockRestrictions' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockCommonShares' xlink:label='loc_TreasuryStockCommonShares' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockPreferredShares' xlink:label='loc_TreasuryStockPreferredShares' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockShares' xlink:label='loc_TreasuryStockShares' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockCommonValue' xlink:label='loc_TreasuryStockCommonValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockPreferredValue' xlink:label='loc_TreasuryStockPreferredValue' xlink:type='locator' />
    <link:loc xlink:href='../elts/us-gaap-2022.xsd#us-gaap_TreasuryStockValue' xlink:label='loc_TreasuryStockValue' xlink:type='locator' />
    <link:presentationArc order='30' preferredLabel='http://www.xbrl.org/2003/role/totalLabel' xlink:arcrole='http://www.xbrl.org/2003/arcrole/parent-child' xlink:from='loc_AccountsPayableAndAccruedLiabilitiesCurrentAbstract' xlink:to='loc_AccountsPayableAndAccruedLiabilitiesCurrent' xlink:type='arc' />
    <link:presentationArc order='20' xlink:arcrole='http://www.xbrl.org/2003/arcrole/parent-child' xlink:from='loc_AccountsPayableAndAccruedLiabilitiesCurrentAbstract' xlink:to='loc_AccruedLiabilitiesCurrentAbstract' xlink:type='arc' />
    <link:presentationArc order='10' xlink:arcrole='http://www.xbrl.org/2003/arcrole/parent-child' xlink:from='loc_AccountsPayableAndAccruedLiabilitiesCurrentAbstract' xlink:to='loc_AccountsPayableCurrent' xlink:type='arc' />
  </link:presentationLink>`

	var presentationLinkBase XLinkPresentationExtendedLink
	require.NoError(t, xml.Unmarshal([]byte(presLink), &presentationLinkBase))

	presentationLinkBase.Role = "test"
	taxonomy := XBRLTaxonomy{AppInfo: XBRLAppInfo{LinkBases: []*XLinkLinkBase{&XLinkLinkBase{PresentationLinks: []XLinkPresentationExtendedLink{presentationLinkBase}}}}}

	require.NoError(t, FormatLegendForPresentationLink(&taxonomy, "test"))
}

type PresentationNode struct {
	Locator        *XLinkLocator
	PreferredLabel string
	Order          int

	Children []*PresentationNode
}

func FormatLegendForPresentationLink(t *XBRLTaxonomy, presentationRole string) error {
	var presentationLink *XLinkPresentationExtendedLink
	var presentationLinkBase *XLinkLinkBase

	// First find the presentation link with the given role
	for _, linkbase := range t.AppInfo.LinkBases {
		for _, pl := range linkbase.PresentationLinks {
			if pl.Role == presentationRole {
				presentationLink = &pl
				presentationLinkBase = linkbase
				break
			}
		}
	}

	if presentationLink == nil {
		return fmt.Errorf("no presentation link with role '%s' was found", presentationRole)
	}

	parentToChildren := map[string][]string{}

	for _, arc := range presentationLink.PresentationArcs {
		parentToChildren[arc.From] = append(parentToChildren[arc.From], arc.To)
	}

	//presentationRoot := &PresentationNode{Name: "(root)"}
	allNodes := make([]*PresentationNode, 0, len(presentationLink.PresentationArcs))

	for _, arc := range presentationLink.PresentationArcs {
		if arc.ArcRole != "http://www.xbrl.org/2003/arcrole/parent-child" {
			return fmt.Errorf("unknown arcrole: %s", arc.ArcRole)
		}

		child := &PresentationNode{Order: arc.Order, PreferredLabel: arc.PreferredLabel}
		parent := &PresentationNode{Children: []*PresentationNode{child}}

		for _, loc := range presentationLink.Locators {
			loc := loc

			if loc.Label == arc.From {
				parent.Locator = &loc
			}

			if loc.Label == arc.To {
				child.Locator = &loc
			}
		}

		if parent.Locator == nil || child.Locator == nil {
			return fmt.Errorf("TODO")
		}

		existingParentNode := findNodeInPresentationMultiGraph(allNodes, parent.Locator.Label)
		if existingParentNode == nil {
			allNodes = append(allNodes, parent)
		} else {
			existingParentNode.Children = append(existingParentNode.Children, child)
		}
	}
	//
	//// Combine nodes?
	//root := &PresentationNode{}
	//
	//for _, node := range allNodes {
	//	// Try to find parent in tree
	//
	//	existingParentInTree := findNodeInPresentationTree(root, node.Locator.Label)
	//	if existingParentInTree != nil {
	//		existingParentInTree.Children = append(existingParentInTree.Children, existingParentInTree.Children...)
	//		if existingParentInTree.Order == "" {
	//			existingParentInTree.Order = node.Order
	//		}
	//
	//	} else {
	//		root.Children = append(root.Children, node)
	//	}
	//}

	// Print I guess
	for _, root := range allNodes {
		printPresentationTree(t, presentationLinkBase.OriginalFile, root, "")
		fmt.Println("---")
	}

	return nil
}

func findNodeInPresentationMultiGraph(roots []*PresentationNode, label string) *PresentationNode {
	for _, root := range roots {
		if node := findNodeInPresentationTree(root, label); node != nil {
			return node
		}
	}

	return nil
}

func findNodeInPresentationTree(root *PresentationNode, label string) *PresentationNode {
	if root.Locator != nil && root.Locator.Label == label {
		return root
	}

	for _, child := range root.Children {
		if node := findNodeInPresentationTree(child, label); node != nil {
			return node
		}
	}

	return nil
}

func printPresentationTree(t *XBRLTaxonomy, rootFile string, root *PresentationNode, prefix string) {
	sort.Slice(root.Children, func(i, j int) bool {
		return root.Children[i].Order < root.Children[j].Order
	})

	if root.Locator != nil {
		labelStr := root.Locator.Label

		element, err := t.ResolveLocator(rootFile, *root.Locator)
		if err == nil {
			labels, err := t.LookupLabels(element.Name)
			if err == nil {
				preferredLabel := root.PreferredLabel
				standardLabel := false
				if preferredLabel == "" {
					preferredLabel = "http://www.xbrl.org/2003/role/label"
					standardLabel = true
				}
				for _, label := range labels {
					if label.Role == preferredLabel || (label.Role == "" && standardLabel) {
						labelStr = label.Value
					}
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}

		fmt.Printf("%s%d: %s\n", prefix, root.Order, labelStr)
	}

	for _, child := range root.Children {
		printPresentationTree(t, rootFile, child, prefix+"  ")
	}
}
