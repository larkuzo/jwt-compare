import csv

import matplotlib.pyplot as plt

COUNT = 50
rsa = []
ecdsa = []
hmac = []
which = [rsa, ecdsa, hmac]

with open('result.csv', newline='') as f:
    reader = csv.reader(f, delimiter=',', quotechar='"')
    index = 0
    for row in reader:
        if len(row) == 0:
            index += 1
            continue

        which[index].append([int(row[0]), int(row[1]), int(row[2])])

fig, ax = plt.subplots()
fig.suptitle('GENERATE TOKEN', fontsize=20)
ax.plot([i for i in range(COUNT)], [rsa[i][0] / 1000000 for i in range(COUNT)], label='RSA')
ax.plot([i for i in range(COUNT)], [ecdsa[i][0] / 1000000 for i in range(COUNT)], label='ECDSA')
ax.plot([i for i in range(COUNT)], [hmac[i][0] / 1000000 for i in range(COUNT)], label='HMAC')
ax.grid()
ax.set(xlabel='PERCOBAAN', ylabel='WAKTU (ms)')
ax.legend()
fig.savefig('generate.png')

fig, ax = plt.subplots()
fig.suptitle('TOKEN SIZE', fontsize=20)
ax.plot([i for i in range(COUNT)], [rsa[i][1] for i in range(COUNT)], label='RSA')
ax.plot([i for i in range(COUNT)], [ecdsa[i][1] for i in range(COUNT)], label='ECDSA')
ax.plot([i for i in range(COUNT)], [hmac[i][1] for i in range(COUNT)], label='HMAC')
ax.grid()
ax.set(xlabel='PERCOBAAN', ylabel='SIZE (bytes)')
ax.legend()
fig.savefig('size.png')

fig, ax = plt.subplots()
fig.suptitle('DATA TRANSFER', fontsize=20)
ax.plot([i for i in range(COUNT)], [rsa[i][2] / 1000000 for i in range(COUNT)], label='RSA')
ax.plot([i for i in range(COUNT)], [ecdsa[i][2] / 1000000 for i in range(COUNT)], label='ECDSA')
ax.plot([i for i in range(COUNT)], [hmac[i][2] / 1000000 for i in range(COUNT)], label='HMAC')
ax.grid()
ax.set(xlabel='PERCOBAAN', ylabel='WAKTU (ms)')
ax.legend()
fig.savefig('transfer.png')

print('DONE')
