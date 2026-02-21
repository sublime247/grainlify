#![cfg(test)]
/// # Escrow Analytics & Monitoring View Tests
///
/// Closes #391
///
/// This module validates that every monitoring metric and analytics view correctly
/// reflects the escrow state after lock, release, and refund operations — including
/// both success and failure/error paths.
///
/// ## Coverage
/// * `get_aggregate_stats`  – totals update after lock → release → refund lifecycle
/// * `get_escrow_count`     – increments on each lock; never decrements
/// * `query_escrows_by_status` – returns correct subset filtered by status
/// * `query_escrows_by_amount` – range filter works for locked, released, and mixed states
/// * `query_escrows_by_deadline` – deadline range filter returns correct bounties
/// * `query_escrows_by_depositor` – per-depositor index is populated on lock
/// * `get_escrow_ids_by_status` – ID-only view mirrors full-object equivalent
/// * `get_refund_eligibility` – eligibility flags flip correctly across lifecycle
/// * `get_refund_history`    – history vector is populated by approved-refund path
/// * Monitoring event emission – lock/release/refund each emit ≥ 1 event
/// * Error flows             – failed attempts do not corrupt metrics
use crate::{BountyEscrowContract, BountyEscrowContractClient, EscrowStatus, RefundMode};
use soroban_sdk::{
    testutils::{Address as _, Events, Ledger},
    token, Address, Env,
};

// ---------------------------------------------------------------------------
// Shared helpers – matching the pattern used in the existing test.rs
// ---------------------------------------------------------------------------

fn create_token_contract<'a>(
    e: &'a Env,
    admin: &Address,
) -> (token::Client<'a>, token::StellarAssetClient<'a>) {
    let contract_address = e.register_stellar_asset_contract(admin.clone());
    (
        token::Client::new(e, &contract_address),
        token::StellarAssetClient::new(e, &contract_address),
    )
}

fn create_escrow_contract<'a>(e: &'a Env) -> BountyEscrowContractClient<'a> {
    let contract_id = e.register_contract(None, BountyEscrowContract);
    BountyEscrowContractClient::new(e, &contract_id)
}

