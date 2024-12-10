# LASFS-elections


This project is to implement the election voting process used by the [LASFS](https://lasfs.org) / Los Angeles Science Fantasy Society.

The LASFS Voting process is similiar to [Ranked Voting](https://en.wikipedia.org/wiki/Ranked_voting), more specifically Ranked-Choice Voting (RCV)/[Instant-runoff Voting](https://en.wikipedia.org/wiki/Instant-runoff_voting)

The voting is usually conducted during a LASFS meeting and depends on the number of positions up for vote, the number of nominees, and number of voting members.

A minimum number of voting members (quorum) is needed to conduct a LASFS Election.

At the start of the LASFS Election, the majority number value is computed from total of valid voting members.

Each position has a list of nominees.  

When the number of nominees equals the number of positions, voting is over, else, a voting period is opened.

Each voting member creates a ballot for the voting position.  (The ballot is a 3-inch by 5-inch colored index card where the color is defined for a specific position).  The ballot is a list of nominees in order of the voting member's preference for that specific position.

The voting period is closed, and ballots are collected and processed in voting rounds.

For each voting round:

The first (1st) valid nominee is read from the ballot, and the nominee vote count is incremented.  The ballots are then collected in piles based on nominee/dead status.

Invalid nominees are striked from the ballot until a valid nominee is found.  If no nominee are found on the ballot, the ballot is considered a dead ballot.  

As soon as the vote count for a nominee matches the majority number value, the nominee wins that position and is removed from the voting process.

The ballots continue to be processed if there are additional nominess and positions.

At the end of the voting round, if there are still positions, the nominee with the smallest voting count is removed from the voting process, and the collection of that nominee's ballots is re-processed in the next voting round.

The voting rounds repeat until all positions are filled.

