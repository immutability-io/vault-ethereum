// Copyright (C) Immutability, LLC - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by Jeff Ploughman <jeff@immutability.io>, August 2019

package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/immutability-io/vault-ethereum/contracts/erc721"
	"github.com/immutability-io/vault-ethereum/util"
)

const erc721Contract string = "erc-721"

// interface ERC721
// {
//   /**
//    * @dev Transfers the ownership of an NFT from one address to another address.
//    * @notice Throws unless `msg.sender` is the current owner, an authorized operator, or the
//    * approved address for this NFT. Throws if `_from` is not the current owner. Throws if `_to` is
//    * the zero address. Throws if `_tokenId` is not a valid NFT. When transfer is complete, this
//    * function checks if `_to` is a smart contract (code size > 0). If so, it calls
//    * `onERC721Received` on `_to` and throws if the return value is not
//    * `bytes4(keccak256("onERC721Received(address,uint256,bytes)"))`.
//    * @param _from The current owner of the NFT.
//    * @param _to The new owner.
//    * @param _tokenId The NFT to transfer.
//    * @param _data Additional data with no specified format, sent in call to `_to`.
//    */
//   function safeTransferFrom(
//     address _from,
//     address _to,
//     uint256 _tokenId,
//     bytes calldata _data
//   )

//   /**
//    * @dev Set or reaffirm the approved address for an NFT.
//    * @notice The zero address indicates there is no approved address. Throws unless `msg.sender` is
//    * the current NFT owner, or an authorized operator of the current owner.
//    * @param _approved The new approved NFT controller.
//    * @param _tokenId The NFT to approve.
//    */
//   function approve(
//     address _approved,
//     uint256 _tokenId
//   )
//     external;

//   /**
//    * @dev Enables or disables approval for a third party ("operator") to manage all of
//    * `msg.sender`'s assets. It also emits the ApprovalForAll event.
//    * @notice The contract MUST allow multiple operators per owner.
//    * @param _operator Address to add to the set of authorized operators.
//    * @param _approved True if the operators is approved, false to revoke approval.
//    */
//   function setApprovalForAll(
//     address _operator,
//     bool _approved
//   )
//     external;

//   /**
//    * @dev Returns the number of NFTs owned by `_owner`. NFTs assigned to the zero address are
//    * considered invalid, and this function throws for queries about the zero address.
//    * @param _owner Address for whom to query the balance.
//    * @return Balance of _owner.
//    */
//   function balanceOf(
//     address _owner
//   )
//     external
//     view
//     returns (uint256);

//   /**
//    * @dev Returns the address of the owner of the NFT. NFTs assigned to zero address are considered
//    * invalid, and queries about them do throw.
//    * @param _tokenId The identifier for an NFT.
//    * @return Address of _tokenId owner.
//    */
//   function ownerOf(
//     uint256 _tokenId
//   )
//     external
//     view
//     returns (address);

//   /**
//    * @dev Get the approved address for a single NFT.
//    * @notice Throws if `_tokenId` is not a valid NFT.
//    * @param _tokenId The NFT to find the approved address for.
//    * @return Address that _tokenId is approved for.
//    */
//   function getApproved(
//     uint256 _tokenId
//   )
//     external
//     view
//     returns (address);

//   /**
//    * @dev Returns true if `_operator` is an approved operator for `_owner`, false otherwise.
//    * @param _owner The address that owns the NFTs.
//    * @param _operator The address that acts on behalf of the owner.
//    * @return True if approved for all, false otherwise.
//    */
//   function isApprovedForAll(
//     address _owner,
//     address _operator
//   )
//     external
//     view
//     returns (bool);

// }

// ERC721Paths are the path handlers for Non-Fungible Tokens
func ERC721Paths(b *PluginBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern:      ContractPath(erc721Contract, "safeTransferFrom"),
			HelpSynopsis: "Transfers the ownership of an NFT from one address to another address.",
			HelpDescription: `

Transfers the ownership of an NFT from one address to another address.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"to": {
					Type:        framework.TypeString,
					Description: "The address to transfer the NFT to.",
				},
				"token_id": {
					Type:        framework.TypeString,
					Description: "The NFT that got transfered.",
				},
				"data": {
					Type:        framework.TypeString,
					Description: "Additional data with no specified format, sent in call to \"to\".",
				},
				"encoding": {
					Type:        framework.TypeString,
					Default:     "utf8",
					Description: "The encoding of the data.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721SafeTransferFrom,
				logical.UpdateOperation: b.pathERC721SafeTransferFrom,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "approve"),
			HelpSynopsis: "Set or reaffirm the approved address for an NFT",
			HelpDescription: `

Set or reaffirm the approved address for an NFT.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"approved": {
					Type:        framework.TypeString,
					Description: "The address to approve as operator of the NFT. Defaults to RLP empty byte sequence..",
					Default:     util.ZeroAddress,
				},
				"token_id": {
					Type:        framework.TypeString,
					Description: "The NFT to approve.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721Approve,
				logical.UpdateOperation: b.pathERC721Approve,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "setApprovalForAll"),
			HelpSynopsis: "Enables or disables approval for a third party",
			HelpDescription: `

Enables or disables approval for a third party ("operator") to manage all of
msg.senders assets.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"operator": {
					Type:        framework.TypeString,
					Description: "Address to add to the set of authorized operators.",
				},
				"approved": {
					Type:        framework.TypeBool,
					Description: "True if the operators is approved, false to revoke approval.",
					Default:     false,
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721SetApprovalForAll,
				logical.UpdateOperation: b.pathERC721SetApprovalForAll,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "balanceOf"),
			HelpSynopsis: "Returns the number of NFTs owned",
			HelpDescription: `

Returns the number of NFTs owned. NFTs assigned to the zero address are
considered invalid, and this function throws for queries about the zero address.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"owner": {
					Type:        framework.TypeString,
					Description: "Address to add to the set of authorized operators.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721BalanceOf,
				logical.UpdateOperation: b.pathERC721BalanceOf,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "ownerOf"),
			HelpSynopsis: "Returns the address of the owner of the NFT.",
			HelpDescription: `

Returns the address of the owner of the NFT. NFTs assigned to zero address are considered
invalid, and queries about them do throw.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"token_id": {
					Type:        framework.TypeString,
					Description: "The identifier for an NFT.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721OwnerOf,
				logical.UpdateOperation: b.pathERC721OwnerOf,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "getApproved"),
			HelpSynopsis: "Get the approved address for a single NFT.",
			HelpDescription: `

Get the approved address for a single NFT.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"token_id": {
					Type:        framework.TypeString,
					Description: "The identifier for an NFT.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721GetApproved,
				logical.UpdateOperation: b.pathERC721GetApproved,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "isApprovedForAll"),
			HelpSynopsis: "Returns true if operator is an approved operator for owner, false otherwise.",
			HelpDescription: `

Returns true if operator is an approved operator for owner, false otherwise.

`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"owner": {
					Type:        framework.TypeString,
					Description: "The address that owns the NFTs",
				},
				"operator": {
					Type:        framework.TypeString,
					Description: "The address that acts on behalf of the owner.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721IsApprovedForAll,
				logical.UpdateOperation: b.pathERC721IsApprovedForAll,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "TokenByIndex"),
			HelpSynopsis: "TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.",
			HelpDescription: `
		
TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
		
		`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"index": {
					Type:        framework.TypeString,
					Description: "The index",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721TokenByIndex,
				logical.UpdateOperation: b.pathERC721TokenByIndex,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "TokenOfOwnerByIndex"),
			HelpSynopsis: "TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.",
			HelpDescription: `
		
		TokenByIndex is a free data retrieval call binding the contract method 0x2f745c59.
		
		`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"owner": {
					Type:        framework.TypeString,
					Description: "Address to add to the set of authorized operators.",
				},
				"index": {
					Type:        framework.TypeString,
					Description: "The index",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721TokenOfOwnerByIndex,
				logical.UpdateOperation: b.pathERC721TokenOfOwnerByIndex,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "Metadata"),
			HelpSynopsis: "This returns TotalSupply, Name, and Symbol.",
			HelpDescription: `
		
		Metadata calls free data retrieval methods at 0x18160ddd, 0x06fdde03, 0x95d89b41.
		This returns TotalSupply, Name, and Symbol 
		
		`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721Metadata,
				logical.UpdateOperation: b.pathERC721Metadata,
			},
		},
		{
			Pattern:      ContractPath(erc721Contract, "TokenURI"),
			HelpSynopsis: "TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.",
			HelpDescription: `
		
		TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
		
		`,
			Fields: map[string]*framework.FieldSchema{
				"name": {Type: framework.TypeString},
				"contract": {
					Type:        framework.TypeString,
					Description: "The address of the ERC-721 NFT.",
				},
				"token_id": {
					Type:        framework.TypeString,
					Description: "The NFT.",
				},
			},
			ExistenceCheck: pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathERC721TokenURI,
				logical.UpdateOperation: b.pathERC721TokenURI,
			},
		},
	}
}

//   /**
//    * @dev Transfers the ownership of an NFT from one address to another address.
//    * @notice Throws unless `msg.sender` is the current owner, an authorized operator, or the
//    * approved address for this NFT. Throws if `_from` is not the current owner. Throws if `_to` is
//    * the zero address. Throws if `_tokenId` is not a valid NFT. When transfer is complete, this
//    * function checks if `_to` is a smart contract (code size > 0). If so, it calls
//    * `onERC721Received` on `_to` and throws if the return value is not
//    * `bytes4(keccak256("onERC721Received(address,uint256,bytes)"))`.
//    * @param _from The current owner of the NFT.
//    * @param _to The new owner.
//    * @param _tokenId The NFT to transfer.
//    * @param _data Additional data with no specified format, sent in call to `_to`.
//    */
//   function safeTransferFrom(
//     address _from,
//     address _to,
//     uint256 _tokenId,
//     bytes calldata _data
//   )

func (b *PluginBackend) pathERC721SafeTransferFrom(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var additionalData []byte
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	tokenID := util.ValidNumber(data.Get("token_id").(string))
	dataOrFile := data.Get("data").(string)
	encoding := data.Get("encoding").(string)
	if encoding == "hex" {
		additionalData, err = util.Decode([]byte(dataOrFile))
		if err != nil {
			return nil, err
		}
	} else if encoding == "utf8" {
		additionalData = []byte(dataOrFile)
	} else {
		return nil, fmt.Errorf("invalid encoding encountered - %s", encoding)
	}

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactionParams, err := b.getBaseData(client, account.Address, data, "to")
	if err != nil {
		return nil, err
	}

	err = config.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	err = accountJSON.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	tx, err := tokenSession.SafeTransferFrom0(account.Address, *transactionParams.Address, tokenID, additionalData)
	if err != nil {
		return nil, err
	}

	var signedTxBuff bytes.Buffer
	tx.EncodeRLP(&signedTxBuff)
	return &logical.Response{
		Data: map[string]interface{}{
			"contract":           tokenAddress.Hex(),
			"token_id":           tokenID.String(),
			"data":               dataOrFile,
			"encoding":           encoding,
			"transaction_hash":   tx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"from":               account.Address.Hex(),
			"to":                 transactionParams.Address.String(),
			"nonce":              tx.Nonce(),
			"gas_price":          tx.GasPrice(),
			"gas_limit":          tx.Gas(),
		},
	}, nil
}

//   /**
//    * @dev Set or reaffirm the approved address for an NFT.
//    * @notice The zero address indicates there is no approved address. Throws unless `msg.sender` is
//    * the current NFT owner, or an authorized operator of the current owner.
//    * @param _approved The new approved NFT controller.
//    * @param _tokenId The NFT to approve.
//    */
//   function approve(
//     address _approved,
//     uint256 _tokenId
//   )

func (b *PluginBackend) pathERC721Approve(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	tokenID := util.ValidNumber(data.Get("token_id").(string))

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactionParams, err := b.getBaseData(client, account.Address, data, "approved")
	if err != nil {
		return nil, err
	}

	err = config.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	err = accountJSON.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	tx, err := tokenSession.Approve(*transactionParams.Address, tokenID)
	if err != nil {
		return nil, err
	}

	var signedTxBuff bytes.Buffer
	tx.EncodeRLP(&signedTxBuff)
	return &logical.Response{
		Data: map[string]interface{}{
			"contract":           tokenAddress.Hex(),
			"token_id":           tokenID.String(),
			"transaction_hash":   tx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"from":               account.Address.Hex(),
			"to":                 transactionParams.Address.String(),
			"nonce":              tx.Nonce(),
			"gas_price":          tx.GasPrice(),
			"gas_limit":          tx.Gas(),
		},
	}, nil
}

//   /**
//    * @dev Enables or disables approval for a third party ("operator") to manage all of
//    * `msg.sender`'s assets. It also emits the ApprovalForAll event.
//    * @notice The contract MUST allow multiple operators per owner.
//    * @param _operator Address to add to the set of authorized operators.
//    * @param _approved True if the operators is approved, false to revoke approval.
//    */
//   function setApprovalForAll(
//     address _operator,
//     bool _approved
//   )
func (b *PluginBackend) pathERC721SetApprovalForAll(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	approved := data.Get("approved").(bool)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactionParams, err := b.getBaseData(client, account.Address, data, "operator")
	if err != nil {
		return nil, err
	}

	err = config.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	err = accountJSON.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	tx, err := tokenSession.SetApprovalForAll(*transactionParams.Address, approved)
	if err != nil {
		return nil, err
	}

	var signedTxBuff bytes.Buffer
	tx.EncodeRLP(&signedTxBuff)
	return &logical.Response{
		Data: map[string]interface{}{
			"contract":           tokenAddress.Hex(),
			"approved":           approved,
			"transaction_hash":   tx.Hash().Hex(),
			"signed_transaction": hexutil.Encode(signedTxBuff.Bytes()),
			"from":               account.Address.Hex(),
			"operator":           transactionParams.Address.String(),
			"nonce":              tx.Nonce(),
			"gas_price":          tx.GasPrice(),
			"gas_limit":          tx.Gas(),
		},
	}, nil
}

//   /**
//    * @dev Returns the number of NFTs owned by `_owner`. NFTs assigned to the zero address are
//    * considered invalid, and this function throws for queries about the zero address.
//    * @param _owner Address for whom to query the balance.
//    * @return Balance of _owner.
//    */
//   function balanceOf(
//     address _owner
//   )

func (b *PluginBackend) pathERC721BalanceOf(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactionParams, err := b.getBaseData(client, account.Address, data, "owner")
	if err != nil {
		return nil, err
	}

	err = config.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	err = accountJSON.ValidAddress(transactionParams.Address)
	if err != nil {
		return nil, err
	}
	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	balance, err := tokenSession.BalanceOf(*transactionParams.Address)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"balance":  balance.String(),
			"owner":    transactionParams.Address.String(),
		},
	}, nil
}

//   /**
//    * @dev Returns the address of the owner of the NFT. NFTs assigned to zero address are considered
//    * invalid, and queries about them do throw.
//    * @param _tokenId The identifier for an NFT.
//    * @return Address of _tokenId owner.
//    */
//   function ownerOf(
//     uint256 _tokenId
//   )

func (b *PluginBackend) pathERC721OwnerOf(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	tokenID := util.ValidNumber(data.Get("token_id").(string))

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	owner, err := tokenSession.OwnerOf(tokenID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"owner":    owner.Hex(),
			"token_id": tokenID.String(),
		},
	}, nil
}

//   /**
//    * @dev Get the approved address for a single NFT.
//    * @notice Throws if `_tokenId` is not a valid NFT.
//    * @param _tokenId The NFT to find the approved address for.
//    * @return Address that _tokenId is approved for.
//    */
//   function getApproved(
//     uint256 _tokenId
//   )

func (b *PluginBackend) pathERC721GetApproved(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	tokenID := util.ValidNumber(data.Get("token_id").(string))

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	approved, err := tokenSession.GetApproved(tokenID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"approved": approved.Hex(),
			"token_id": tokenID.String(),
		},
	}, nil
}

//   /**
//    * @dev Returns true if `_operator` is an approved operator for `_owner`, false otherwise.
//    * @param _owner The address that owns the NFTs.
//    * @param _operator The address that acts on behalf of the owner.
//    * @return True if approved for all, false otherwise.
//    */
//   function isApprovedForAll(
//     address _owner,
//     address _operator
//   )

func (b *PluginBackend) pathERC721IsApprovedForAll(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}
	transactionParams, err := b.getBaseData(client, account.Address, data, "owner")
	if err != nil {
		return nil, err
	}
	owner := *transactionParams.Address
	transactionParams, err = b.getBaseData(client, account.Address, data, "operator")
	if err != nil {
		return nil, err
	}
	operator := *transactionParams.Address

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	approved, err := tokenSession.IsApprovedForAll(owner, operator)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"approved": approved,
			"operator": operator.Hex(),
			"owner":    owner.Hex(),
		},
	}, nil
}

func (b *PluginBackend) pathERC721TokenByIndex(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}
	index := util.ValidNumber(data.Get("index").(string))

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	tokenID, err := tokenSession.TokenByIndex(index)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"token":    tokenID.String(),
			"index":    index.String(),
		},
	}, nil
}

func (b *PluginBackend) pathERC721TokenOfOwnerByIndex(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}
	index := util.ValidNumber(data.Get("index").(string))

	transactionParams, err := b.getBaseData(client, account.Address, data, "owner")
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	tokenID, err := tokenSession.TokenOfOwnerByIndex(*transactionParams.Address, index)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"owner":    transactionParams.Address.Hex(),
			"index":    index.String(),
			"token":    tokenID.String(),
		},
	}, nil
}

func (b *PluginBackend) pathERC721Metadata(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	supply, err := tokenSession.TotalSupply()
	if err != nil {
		return nil, err
	}
	nftName, err := tokenSession.Name()
	if err != nil {
		return nil, err
	}
	symbol, err := tokenSession.Symbol()
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"contract": tokenAddress.Hex(),
			"name":     nftName,
			"symbol":   symbol,
			"supply":   supply.String(),
		},
	}, nil
}
func (b *PluginBackend) pathERC721TokenURI(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}
	name := data.Get("name").(string)
	tokenID := util.ValidNumber(data.Get("token_id").(string))

	accountJSON, err := readAccount(ctx, req, name)
	if err != nil {
		return nil, err
	}
	wallet, account, err := getWalletAndAccount(*accountJSON)
	if err != nil {
		return nil, err
	}

	tokenAddress := common.HexToAddress(data.Get("contract").(string))

	chainID := util.ValidNumber(config.ChainID)
	if chainID == nil {
		return nil, fmt.Errorf("invalid chain ID")
	}

	client, err := ethclient.Dial(config.getRPCURL())
	if err != nil {
		return nil, err
	}

	instance, err := erc721.NewErc721(tokenAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{}

	transactOpts, err := b.NewWalletTransactor(chainID, wallet, account)
	if err != nil {
		return nil, err
	}

	//transactOpts needs gas etc.
	tokenSession := &erc721.Erc721Session{
		Contract:     instance,  // Generic contract caller binding to set the session for
		CallOpts:     *callOpts, // Call options to use throughout this session
		TransactOpts: *transactOpts,
	}

	tokenURI, err := tokenSession.TokenURI(tokenID)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"contract":  tokenAddress.Hex(),
			"token_uri": tokenURI,
		},
	}, nil
}
