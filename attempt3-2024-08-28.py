# https://replit.com/@RobPowell1/LASFSVotingProcess#attempt3.py
import collections

"""

References:
* https://en.wikipedia.org/wiki/Instant-runoff_voting

"""

class LASFSMemberBallot:
  """LASFSMemberBallot

  A ballot is 'dead' when it has no further nominees

  """
  def __init__(self, name, nominees=[]):
    self.name = name
    self.original_nominees = nominees
    self.nominees = collections.OrderedDict([[n, True] for n in nominees])
    self.dead=False

  def __repr__(self):
    r"""String representation

    String representation constructor test

    >>> LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"]).__repr__()
    "name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])"


    String representation using injected data

    >>> LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])


    >>> l1 = LASFSMemberBallot("Elayne", [])

    >>> l1
    name:Elayne votes:OrderedDict()


    >>> l2 = LASFSMemberBallot("Elayne", [])

    >>> l2.nominees = collections.OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])

    >>> l2
    name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])


    Injecting votes with struck nominees

    >>> l3 = LASFSMemberBallot("Elayne", [])

    >>> l3.nominees = collections.OrderedDict([('Matthew', False), ('Michael', False), ('John', True)])

    >>> l3
    name:Elayne votes:OrderedDict([('Matthew', False), ('Michael', False), ('John', True)])

    """

    return "name:{} votes:{}".format(self.name, self.nominees)

  def isValid(self):
    """

    >>> b1 = LASFSMemberBallot("emptyballot",[])
    >>> b1.isValid()
    False

    """
    if (
      (len(self.name) != 0) and
      (self.nominees != None) and
      (len(self.nominees) != 0)
    ):
      return True
    else:
      return False


  def getPreferredNominee(self):
    """getPreferredNominee

    >>> e1 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e1.getPreferredNominee()
    'Matthew'


    >>> e2 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e2.nominees["Matthew"] = False
    >>> e2.getPreferredNominee()
    'Michael'


    >>> e3 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e3.nominees = collections.OrderedDict([('Matthew', False), ('Michael', False), ('John', True)])
    >>> e3.getPreferredNominee()
    'John'


    >>> deadballot1 = LASFSMemberBallot("deadballottest",[])
    >>> deadballot1.dead
    False
    >>> deadballot1.getPreferredNominee()

    >>> deadballot1.dead
    True


    """

    if not self.dead:
      try:
        return [n for n in self.nominees if self.nominees[n]][0]
      except IndexError:
        self.dead = True
        return None

  def strikeNomineeByName(self, nominee):
    """strikeNominee(nominee)

    Create a LASFSMemberBallot and strike off names
    >>> e3 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> for i in range(4):
    ...   print("ballot is dead: {}".format(e3.dead))
    ...   n = e3.getPreferredNominee()
    ...   print(n)
    ...   e3.strikeNomineeByName(n)
    ballot is dead: False
    Matthew
    ballot is dead: False
    Michael
    ballot is dead: False
    John
    ballot is dead: True
    None


    Scenario 1

    >>> e1 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e1.strikeNomineeByName("Matthew")
    >>> e1
    name:Elayne votes:OrderedDict([('Matthew', False), ('Michael', True), ('John', True)])

    Scenario 2

    >>> e2 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e2.strikeNomineeByName("Matthew")
    >>> e2.strikeNomineeByName("Michael")
    >>> e2
    name:Elayne votes:OrderedDict([('Matthew', False), ('Michael', False), ('John', True)])

    Scenario 3

    Striking all valid names on a ballot makes the ballot dead.

    >>> e3 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e3.strikeNomineeByName("Matthew")
    >>> e3.strikeNomineeByName("Michael")
    >>> e3.dead
    False
    >>> e3.strikeNomineeByName("John")
    >>> e3.dead
    True

    Scenario 4

    Striking an invalid name does not change ballot state.

    >>> e4 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e4
    name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])
    >>> e4.dead
    False
    >>> e4.strikeNomineeByName("Larry")
    >>> e4
    name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])
    >>> e4.dead
    False


    """

    if nominee in self.nominees:
      self.nominees[nominee]=False

    if (not self.dead) and (True not in self.nominees.values()):
      self.dead = True

    #self.nominees.pop(nominee, None)

  def strikePreferredNominee(self):
    """strikeNominee()

    >>> e = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
    >>> e.strikePreferredNominee()
    >>> e
    name:Elayne votes:OrderedDict([('Matthew', False), ('Michael', True), ('John', True)])

    """
    self.strikeNomineeByName(self.getPreferredNominee())


class LASFSElectionRoundResults:
  """

  >>> r1 = LASFSElectionRoundResults()

  >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"]))

  >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Albert", ["Michael", "Matthew", "John"]))

  >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Matthew", "Michael", "John"]))

  >>> print(r1.summary())
  ['Matthew', 2]
  ['Michael', 1]


  # sort results with hightest count test
  >>> r2 = LASFSElectionRoundResults()

  >>> r2.processLASFSMemberBallot(LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"]))

  >>> r2.processLASFSMemberBallot(LASFSMemberBallot("Albert", ["Michael", "Matthew", "John"]))

  >>> r2.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Michael", "Matthew", "John"]))

  >>> print(r2.summary())
  ['Michael', 2]
  ['Matthew', 1]


  # sort results with hightest count test and then alphabetical
  >>> r3 = LASFSElectionRoundResults()

  >>> r3.processLASFSMemberBallot(LASFSMemberBallot("Elayne",["Matthew", "Michael", "John"]))

  >>> r3.processLASFSMemberBallot(LASFSMemberBallot("Albert",["Michael", "Matthew", "John"]))

  >>> r3.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Michael", "Matthew", "John"]))

  >>> r3.processLASFSMemberBallot(LASFSMemberBallot("June",  ["Matthew", "Michael", "John"]))

  >>> print(r3.summary())
  ['Matthew', 2]
  ['Michael', 2]


  # sort with struck ballot
  >>> r4 = LASFSElectionRoundResults()

  >>> l4 = LASFSMemberBallot("Elayne",[])

  >>> l4.nominees = collections.OrderedDict([('Matthew', False), ('Michael', True), ('John', True)])

  >>> r4.processLASFSMemberBallot(l4)

  >>> r4.processLASFSMemberBallot(LASFSMemberBallot("Albert", ["Michael", "Matthew", "John"]))

  >>> r4.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Matthew", "Michael", "John"]))

  >>> print(r4.summary())
  ['Michael', 2]
  ['Matthew', 1]



  >>> r5 = LASFSElectionRoundResults()

  >>> r5.processLASFSMemberBallot(LASFSMemberBallot("Elayne",   ["Matthew", "Michael", "John"]))

  >>> r5.processLASFSMemberBallot(LASFSMemberBallot("Marty",    ["Matthew", "Michael", "John"]))

  >>> r5.processLASFSMemberBallot(LASFSMemberBallot("June",     ["Michael", "Matthew", "John"]))

  >>> r5.processLASFSMemberBallot(LASFSMemberBallot("Michelle", ["Michael", "Matthew", "John"]))

  >>> r5.processLASFSMemberBallot(LASFSMemberBallot("Albert",   ["John", "Matthew", "Michael"]))

  >>> print(r5.summary())
  ['Matthew', 2]
  ['Michael', 2]
  ['John', 1]


  """

  def __init__(self):
      self.nominees = collections.OrderedDict()

  def processLASFSMemberBallot(self, lasfsMemberBallot):
    """processLASFSMemberBallot

    >>> LASFSElectionRoundResults() # doctest: +ELLIPSIS
    <__main__.LASFSElectionRoundResults ... at ...>


    >>> r1 = LASFSElectionRoundResults()

    >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Elayne",[ "Matthew", "Michael", "John" ]))

    """
    n = lasfsMemberBallot.getPreferredNominee()
    self.nominees[n] = self.nominees.get(n, []) + [lasfsMemberBallot]

    self.nominees = collections.OrderedDict(sorted(self.nominees.items(),key=lambda x: (-len(x[1]), x[0])))


  def strikeAndProcessLASFSMemberBallots(self, nominee):
    """
    >>> s1 = LASFSElectionRoundResults()

    >>> s1.processLASFSMemberBallot(LASFSMemberBallot("Elayne",   ["Matthew", "Michael", "John"]))

    >>> s1.processLASFSMemberBallot(LASFSMemberBallot("Marty",    ["Matthew", "Michael", "John"]))

    >>> s1.processLASFSMemberBallot(LASFSMemberBallot("June",     ["Michael", "Matthew", "John"]))

    >>> s1.processLASFSMemberBallot(LASFSMemberBallot("Michelle", ["Michael", "Matthew", "John"]))

    >>> s1.processLASFSMemberBallot(LASFSMemberBallot("Albert",   ["John", "Michael", "Matthew"]))

    >>> print(s1.summary())
    ['Matthew', 2]
    ['Michael', 2]
    ['John', 1]

    >>> s1.strikeAndProcessLASFSMemberBallots("John")

    >>> print(s1.summary())
    ['Michael', 3]
    ['Matthew', 2]

    """
    ballots = self.nominees.pop(nominee)
    [b.strikeNomineeByName(nominee) for b in ballots]
    for ballot in ballots:
      self.processLASFSMemberBallot(ballot)

  def summary(self):
    return "\n".join([
        str([
                n, len(self.nominees[n])
        ])
        for n in self.nominees.keys()
    ])


class LASFSElectionForPosition:
  def __init__(self, position, nominees):
    self.position = position
    self.nominees = nominees
    self.votingrounds = []
    self.votes = []

  def processLASFSMemberVote(lasfsMemberVote):
    n = lasfsMemberVote.getNextPref()


  def __str__(self):
    r"""String representation

    >>> LASFSElectionForPosition("President", ["Matthew", "Michael", "John"]).__str__()
    "position: President nominees:['Matthew', 'Michael', 'John']\nvotingrounds:[]"

    >>> print(LASFSElectionForPosition("President", ["Matthew", "Michael", "John"]))
    position: President nominees:['Matthew', 'Michael', 'John']
    votingrounds:[]


    """

    return "\n".join((
      "position: {} nominees:{}".format(
        self.position, self.nominees
      ),
      "votingrounds:" + str([ 
        "{}:".format(i) + v for i, v in enumerate(self.votingrounds)
      ])
    ))


if __name__ == '__main__':
  import doctest
  doctest.testmod(verbose=False)
  #doctest.testmod(verbose=True)
  #l = LASFSElectionForPosition("President", ["Matthew", "Michael", "John"])
  #print(l)

