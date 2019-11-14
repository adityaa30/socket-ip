import socket
import time

HEADER_SIZE = 10    

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind((socket.gethostname(), 1235))
s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
s.listen(5)

while True:
    conn, addr = s.accept()
    # print(f'Connected to {addr}')
    print "Connected to" + str(addr)

    msg = 'Welcome to server!'
    # msg = f'{len(msg):<{HEADER_SIZE}}{msg}'
    msg = str(len(msg)).rjust(HEADER_SIZE) + msg

    conn.send(msg.encode('utf-8'))

    while True:
        time.sleep(2)
        # msg = f'Current time : {time.time()}'
        # msg = f'{len(msg):<{HEADER_SIZE}}{msg}'
        msg = 'Current time' + str(time.time())
        msg = str(len(msg)).rjust(HEADER_SIZE) + msg

        conn.send(msg.encode('utf-8'))
