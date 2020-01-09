# PartnerAtStake
Implantation of a new consensus algorithm for blockchain as a semester project
CONSENSUS ALORITHOM         (Proof of Partner or Partner at Stake)
Team:
•	Muhammad Muneeb 16I0128
•	Fahad Mumtaz             16I0098
•	Khizar Shabir                16I0294
Inspiration 
WeChat only approve a new user for joining its platform if the new user know somebody who already uses WeChat. New user has to get the recommendation of early user of WeChat for joining its platform. This was the basic idea that evolved into a consensus algorithm.
Basic Idea
There will be first (n) number of members that will join the blockChain and will be considered old and trusted members of blockChain. After that if any user wants to connect or become a member of blockChain, he will have to ask a previous member for making him his partner. If the old member choose to accept him as his partner, then the new user will be made partner off that old user and the new user will be part of the blockChain. He can make transactions after that.
Consequence of making a new partner
Since a trusted user has allowed a new user to come into the blockChain. Now they are partners and their fate are connected. If a new user does any invalid transaction, it will not be processed, further the one who did the invalid transaction will be given first and final warning. If he does it again. Person who did these invalid transactions and all its partners will be blacklisted and they will not be allowed to process any transaction in blockChain further in the future. That is why it is called “Partner at Stake”.
Choosing Miner
Since a user is putting himself on stake for making a new partner (He will be blacklisted if his partner does any invalid transactions). A user’s partners will be given priority for processing his transaction or mining. User randomly chooses among his partners a miner for his next transaction. That’s how becoming partners is also beneficial. Partners can benefit from each other through coinbase transaction. The name of our coin is FreeCoin, and one coinbase transaction is worth 75 FreeCoins.
Benefits/Advantages
•	Only trusted members will be permitted into the blockChain.
•	Even after entering blockChain, the fate of new member is connected with his partner, to ensure he will not commit fraud.
•	If a user does fraud, he and his members will be blacklisted after one warning.
•	Increasing your partners is also beneficial because you will be preferred over others for your member’s transaction. Hence more mining opportunities.


Further improvements:
•	Convergence of all transactions comes to one node for validation. Because of one coin store. It works in our case but practically the congestion of traffic is possible. We can broadcast coin store on each transaction to all the members and all the members can locally validate their transaction. But inconsistency can be present in the blockChain network. 
•	 Partner is mining transaction of his partners, so we can further increase the security level and we can validate further a mined transaction from non-partner.
