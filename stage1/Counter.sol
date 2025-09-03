// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.0;

contract Counter {
    uint256 public count;

    function increment() external {
        count += 1;
    }
}