export const NATION_CONTRACT_ABI = [
    {
        "constant": false,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "joinNation",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "numNations",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_citizenAddress",
                "type": "address"
            },
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "checkCitizen",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "getNationGovernance",
        "outputs": [
            {
                "name": "_decisionMakingProcess",
                "type": "string"
            },
            {
                "name": "_diplomaticRecognition",
                "type": "bool"
            },
            {
                "name": "_governanceService",
                "type": "string"
            },
            {
                "name": "_nonCitizenUse",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            },
            {
                "name": "_decisionMakingProcess",
                "type": "string"
            },
            {
                "name": "_diplomaticRecognition",
                "type": "bool"
            },
            {
                "name": "_governanceService",
                "type": "string"
            },
            {
                "name": "_nonCitizenUse",
                "type": "bool"
            }
        ],
        "name": "SetNationGovernance",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            },
            {
                "name": "_nationCode",
                "type": "string"
            },
            {
                "name": "_nationCodeLink",
                "type": "string"
            },
            {
                "name": "_lawEnforcementMechanism",
                "type": "string"
            },
            {
                "name": "_profit",
                "type": "bool"
            }
        ],
        "name": "SetNationPolicy",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_nationName",
                "type": "string"
            },
            {
                "name": "_nationDescription",
                "type": "string"
            },
            {
                "name": "_exists",
                "type": "bool"
            },
            {
                "name": "_virtualNation",
                "type": "bool"
            }
        ],
        "name": "createNationCore",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "leaveNation",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_newNation",
                "type": "address"
            }
        ],
        "name": "upgradeNation",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "getNationName",
        "outputs": [
            {
                "name": "",
                "type": "string"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "NationCoreVersion",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "getInitializationBlock",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "owner",
        "outputs": [
            {
                "name": "",
                "type": "address"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "getNationPolicy",
        "outputs": [
            {
                "name": "_nationCode",
                "type": "string"
            },
            {
                "name": "_nationCodeLink",
                "type": "string"
            },
            {
                "name": "_lawEnforcementMechanism",
                "type": "string"
            },
            {
                "name": "_profit",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_newOwner",
                "type": "address"
            }
        ],
        "name": "changeOwner",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "getNationCore",
        "outputs": [
            {
                "name": "_nationName",
                "type": "string"
            },
            {
                "name": "_nationDescription",
                "type": "string"
            },
            {
                "name": "_exists",
                "type": "bool"
            },
            {
                "name": "_virtualNation",
                "type": "bool"
            },
            {
                "name": "_founder",
                "type": "address"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "nationImpl",
        "outputs": [
            {
                "name": "",
                "type": "address"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_founder",
                "type": "address"
            }
        ],
        "name": "getFoundedNations",
        "outputs": [
            {
                "name": "",
                "type": "uint256[]"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "_owner",
                "type": "address"
            }
        ],
        "name": "initialize",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [
            {
                "name": "_nationId",
                "type": "uint256"
            }
        ],
        "name": "getNumCitizens",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "constructor"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "newNation",
                "type": "address"
            }
        ],
        "name": "UpgradeNation",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "newOwner",
                "type": "address"
            }
        ],
        "name": "OwnerChanged",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "founder",
                "type": "address"
            },
            {
                "indexed": false,
                "name": "nationName",
                "type": "string"
            },
            {
                "indexed": true,
                "name": "nationId",
                "type": "uint256"
            }
        ],
        "name": "NationCoreCreated",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "name": "nationName",
                "type": "string"
            },
            {
                "indexed": true,
                "name": "nationId",
                "type": "uint256"
            }
        ],
        "name": "NationPolicySet",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "name": "nationName",
                "type": "string"
            },
            {
                "indexed": true,
                "name": "nationId",
                "type": "uint256"
            }
        ],
        "name": "NationGovernanceSet",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "nationId",
                "type": "uint256"
            },
            {
                "indexed": false,
                "name": "citizenAddress",
                "type": "address"
            }
        ],
        "name": "CitizenJoined",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "name": "nationId",
                "type": "uint256"
            },
            {
                "indexed": false,
                "name": "citizenAddress",
                "type": "address"
            }
        ],
        "name": "CitizenLeft",
        "type": "event"
    }
];
