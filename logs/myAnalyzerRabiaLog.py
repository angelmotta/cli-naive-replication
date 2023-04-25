import csv
import itertools

zip_longest = itertools.zip_longest

f1 = open("./rabia/rabiasvr1log.txt")
f2 = open("./rabia/rabiasvr2log.txt")
f3 = open("./rabia/rabiasvr3log.txt")

csv_f1 = csv.reader(f1, delimiter=" ")
csv_f2 = csv.reader(f2, delimiter=" ")
csv_f3 = csv.reader(f3, delimiter=" ")

countDiff = 0
numLine = 1
setCmds = 0
getCmds = 0
for rowf1, rowf2, rowf3 in zip_longest(csv_f1, csv_f2, csv_f3):
    if rowf1 is None or len(rowf1) < 5:
        # Not SET nor GET could be 'OK' or another command logged
        numLine += 1
        continue

    # Same type of operation in all Rabia Servers
    opR1, opR2, opR3 = rowf1[3], rowf2[3], rowf3[3]
    if opR1 != opR2 or opR1 != opR3 or opR2 != opR3:
        print("Tipo Operacion diferente: Op# " + str(numLine))
        countDiff += 1
        continue
    # Verify operation key in Get Operation
    if len(rowf1) == 5:
        # Get Command
        getCmds += 1
        # Check key of get operation
        keyR1, KeyR2, KeyR3 = rowf1[4], rowf2[4], rowf3[4]
        if keyR1 != KeyR2 or keyR1 != KeyR3 or KeyR2 != KeyR3:
            print("Operacion GET diferente: Op# " + str(numLine))
            countDiff += 1
            continue
    # Verify operation key and value in Set Operation
    elif len(rowf1) == 6:
        # Set Command
        setCmds += 1
        # Check key of set operation
        keyR1, KeyR2, KeyR3 = rowf1[4], rowf2[4], rowf3[4]
        if keyR1 != KeyR2 or keyR1 != KeyR3 or KeyR2 != KeyR3:
            print("Operacion SET diferente key: Op# " + str(numLine))
            countDiff += 1
            continue
        # Check value of set operation
        valR1, valR2, valR3 = rowf1[5], rowf2[5], rowf3[5]
        if valR1 != valR2 or valR1 != valR3 or valR2 != valR3:
            print("Operacion SET diferente value: Op# " + str(numLine))
            countDiff += 1
            continue
    print(rowf1)
    numLine += 1

print("Total Sets Commands:" + str(setCmds))
print("Total Gets Commands:" + str(getCmds))
print("Total Commands: " + str(getCmds + setCmds))
print("Total Inconsistencias: " + str(countDiff))

f1.close()
f2.close()
f3.close()