# https://replit.com/@RobPowell1/LASFSVotingProcess#attempt1.py
import collections
import copy
import json
"""

References:
* https://en.wikipedia.org/wiki/Instant-runoff_voting

"""


class LASFSMemberBallot:
  """LASFSMemberBallot

    A ballot is 'dead' when it has no further nominees

    """

  def __init__(self, name="", nominees=[]):
    """ __init__

          Constructor Test
          ----------------

          >>> ballotconstructor1 = LASFSMemberBallot()

          >>> ballotconstructor1.name
          ''

          >>> ballotconstructor1.nominees
          OrderedDict()

          >>> ballotconstructor1.dead
          True

          Testing empty ballot initialization
          ===================================

          Testing empty ballot initialization
          -----------------------------------

          >>> emptyballot1 = LASFSMemberBallot("emptyballot", [])

          >>> emptyballot1.dead
          True

          Testing valid ballot initialization b1
          --------------------------------------

          >>> b1 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])

          >>> b1.name
          'Elayne'
          >>> b1.nominees
          OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])
          >>> b1.nominees.keys()
          odict_keys(['Matthew', 'Michael', 'John'])
          >>> list(b1.nominees.keys())
          ['Matthew', 'Michael', 'John']
          >>> b1.dead
          False

          Testing valid ballot initialization b2
          --------------------------------------

          >>> b2 = LASFSMemberBallot("Elayne", collections.OrderedDict([('Matthew', True), ('Michael', True), ('John', True)]))

          >>> b2.name
          'Elayne'

          >>> b2.nominees
          OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])
          >>> b2.nominees.keys()
          odict_keys(['Matthew', 'Michael', 'John'])
          >>> b2.dead
          False

          """

    self.name = name
    if isinstance(nominees, list):
      self.nominees = collections.OrderedDict([[n, True] for n in nominees])
      self.dead = not (True in self.nominees.values())

    elif isinstance(nominees, collections.OrderedDict):
      self.nominees = copy.deepcopy(nominees)
      self.dead = not (True in self.nominees.values())

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

        >>> l2 = LASFSMemberBallot("Elayne", collections.OrderedDict([('Matthew', True), ('Michael', True), ('John', True)]))

        >>> l2
        name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])


        Injecting votes with struck nominees
        ------------------------------------

        >>> l3 = LASFSMemberBallot("Elayne", [])

        >>> l3.nominees = collections.OrderedDict([('Matthew', False), ('Michael', False), ('John', True)])

        >>> l3
        name:Elayne votes:OrderedDict([('Matthew', False), ('Michael', False), ('John', True)])

        """

    return "name:{} votes:{}".format(self.name, self.nominees)

  def fromJSON(self, jsonstr):
    """
        >>> j1 = LASFSMemberBallot().fromJSON('{"name":"Elayne", "nominees":["Matthew", "Michael", "John"]}')

        >>> j1
        name:Elayne votes:OrderedDict([('Matthew', True), ('Michael', True), ('John', True)])

        """

    try:
      ballotjson = json.loads(jsonstr)
    except json.JSONDecodeError as e:
      # TODO: logmessage("Invalid JSON syntax:", e)
      raise e

    try:
      name = ballotjson["name"]
    except KeyError as k:
      # TODO: logmessage("Missing 'name' key in JSON data")
      raise k
    assert isinstance(name, str)

    self.name = ballotjson["name"]

    try:
      name = ballotjson["nominees"]
    except KeyError as k:
      # TODO: logmessage("Missing 'nominees' key in JSON data")
      raise k
    assert isinstance(ballotjson["nominees"], list)

    for n in ballotjson["nominees"]:
      assert isinstance(n, str)

    self.nominees = collections.OrderedDict([[n, True]
                                             for n in ballotjson["nominees"]])
    self.dead = not (True in self.nominees.values())
    return self

  def isValid(self):
    """

        >>> b1 = LASFSMemberBallot("emptyballot",[])
        >>> b1.isValid()
        False

        """
    if ((len(self.name) != 0) and (self.nominees != None)
        and (len(self.nominees) != 0)):
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

        >>> e4 = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
        >>> e4.nominees = collections.OrderedDict([('Matthew', False), ('Michael', False), ('John', False)])
        >>> e4.getPreferredNominee()

        >>> deadballot1 = LASFSMemberBallot("deadballottest",[])
        >>> deadballot1.dead
        True

        >>> deadballot1.getPreferredNominee()

        >>> deadballot1.dead
        True


        """

    if not self.dead:
      nominee = None
      try:
        nominee = next((n for n in self.nominees if self.nominees[n]), None)
      except IndexError:
        self.dead = True
        nominee = None
      return nominee
    else:
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
      self.nominees[nominee] = False

    if (not self.dead) and (True not in self.nominees.values()):
      self.dead = True

    # self.nominees.pop(nominee, None)

  def strikePreferredNominee(self):
    """strikeNominee()

        >>> e = LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"])
        >>> e.strikePreferredNominee()
        >>> e
        name:Elayne votes:OrderedDict([('Matthew', False), ('Michael', True), ('John', True)])

        """
    self.strikeNomineeByName(self.getPreferredNominee())


class LASFSElectionRound:
  """

    """

  def __init__(self):
    self.nominees = collections.OrderedDict()
    self.ballots = []

  def summary(self):
    return "\n".join(
        [str([n, len(self.nominees[n])]) for n in self.nominees.keys()])

  def summaryJSON(self):
    return json.dumps(
        [str([n, len(self.nominees[n])]) for n in self.nominees.keys()])

  def processLASFSMemberBallot(self, lasfsMemberBallot):
    """processLASFSMemberBallot

        >>> LASFSElectionRound() # doctest: +ELLIPSIS
        <__main__.LASFSElectionRound ... at ...>

        >>> r1 = LASFSElectionRound()

        >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"]))

        >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Albert", ["Michael", "Matthew", "John"]))

        >>> r1.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Matthew", "Michael", "John"]))

        >>> print(r1.summary())
        ['Matthew', 2]
        ['Michael', 1]

        >>> print(r1.summaryJSON())
        ["[\'Matthew\', 2]", "[\'Michael\', 1]"]

        # sort results with hightest count test
        >>> r2 = LASFSElectionRound()

        >>> r2.processLASFSMemberBallot(LASFSMemberBallot("Elayne", ["Matthew", "Michael", "John"]))

        >>> r2.processLASFSMemberBallot(LASFSMemberBallot("Albert", ["Michael", "Matthew", "John"]))

        >>> r2.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Michael", "Matthew", "John"]))

        >>> print(r2.summary())
        ['Michael', 2]
        ['Matthew', 1]

        >>> print(r2.summaryJSON())
        ["[\'Michael\', 2]", "[\'Matthew\', 1]"]

        # sort results with hightest count test and then alphabetical
        >>> r3 = LASFSElectionRound()

        >>> r3.processLASFSMemberBallot(LASFSMemberBallot("Elayne",["Matthew", "Michael", "John"]))

        >>> r3.processLASFSMemberBallot(LASFSMemberBallot("Albert",["Michael", "Matthew", "John"]))

        >>> r3.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Michael", "Matthew", "John"]))

        >>> r3.processLASFSMemberBallot(LASFSMemberBallot("June",  ["Matthew", "Michael", "John"]))

        >>> print(r3.summary())
        ['Matthew', 2]
        ['Michael', 2]


        # sort with struck ballot
        >>> r4 = LASFSElectionRound()

        >>> l4 = LASFSMemberBallot("Elayne", collections.OrderedDict([('Matthew', False), ('Michael', True), ('John', True)]))

        >>> r4.processLASFSMemberBallot(l4)

        >>> r4.processLASFSMemberBallot(LASFSMemberBallot("Albert", ["Michael", "Matthew", "John"]))

        >>> r4.processLASFSMemberBallot(LASFSMemberBallot("Marty", ["Matthew", "Michael", "John"]))

        >>> print(r4.summary())
        ['Michael', 2]
        ['Matthew', 1]


        >>> r5 = LASFSElectionRound()

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
    n = lasfsMemberBallot.getPreferredNominee()
    self.nominees[n] = self.nominees.get(n, []) + [lasfsMemberBallot]

    self.nominees = collections.OrderedDict(
        sorted(self.nominees.items(), key=lambda x: (-len(x[1]), x[0])))

  def strikeAndProcessLASFSMemberBallots(self, nominee):
    """
        >>> s1 = LASFSElectionRound()

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


class LASFSElectionForPosition:

  def __init__(self, position, nominees):
    self.position = position
    self.nominees = nominees
    self.ballots = []
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

        ---

        >>> e1 = LASFSElectionForPosition("President", ["Matthew", "Michael", "John"])

        >>> e1.__str__()
        "position: President nominees:['Matthew', 'Michael', 'John']\nvotingrounds:[]"

        ---

        """

    return "\n".join([
        "position: {} nominees:{}".format(self.position,
                                          self.nominees), "votingrounds:" +
        str(["{}:".format(i) + v for i, v in enumerate(self.votingrounds)])
    ])


def LASFSVotingDocTest():
  import doctest
  #doctest.testmod(verbose=False)
  doctest.testmod(verbose=True)


if __name__ == '__main__':
  LASFSVotingDocTest()
  # l = LASFSElectionForPosition("President", ["Matthew", "Michael", "John"])
  # print(l)

