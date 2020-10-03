pragma solidity 0.6.0;

library SafeMath {
    /**
     * @dev Returns the multiplication of two unsigned integers, reverting on
     * overflow.
     */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");

        return c;
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * Reverts when dividing by zero.
     */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        return mod(a, b, "SafeMath: modulo by zero");
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * Reverts with custom message when dividing by zero.
     */
    function mod(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b != 0, errorMessage);
        return a % b;
    }
}

contract Paillier {
    using SafeMath for uint256;

    address public owner;

    uint256 public N2 = 1;
    mapping(bytes32 => uint256) results;

    event Add(bytes32 taskId, uint256 result);

    /**
     * @dev constructor
     */
    constructor() public {
        owner = msg.sender;
    }

    /**
     * @dev set N*N, only owner can operate
     */
    function setN2(uint256 _n2) public {
//        require(owner == msg.sender, "only owner");
        N2 = _n2;
    }

    // reset N2 with taskId
    function resetN2(uint256 _n2, bytes32 taskId) public {
        N2 = _n2;
        results[taskId] = 1;
    }

    /**
     * @dev set owner, owner can update N2 params
     */
    function setOwner(address _newOwner) public {
        require(owner == msg.sender, "only owner");
        owner = _newOwner;
    }

    /**
     * @dev returns result by taskId
     */
    function queryResult(bytes32 taskId) view public returns (uint256 result)  {
        return results[taskId];
    }

    /**
    * @dev clear result by taskId
    */
    function clear(bytes32 taskId) public {
//        require(owner == msg.sender, "only owner");
        results[taskId] = 1;
    }

    /**
     * @dev for paillier add
     */
    function paillierAdd(bytes32 taskId, uint256 encryptInput) public {
        if (results[taskId] <= 0) {
            results[taskId] = 1;
        }

        uint256 res = results[taskId];
        res = res.mul(encryptInput);
        results[taskId] = res.mod(N2);
    }
}