import socket

HEADER_SIZE = 10

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((socket.gethostname(), 1235))

while True:

    full_msg = ''
    got_msg_header = False

    while True:
        msg = s.recv(16)

        if not got_msg_header:
            msg_len = int(msg[:HEADER_SIZE].strip())
            got_msg_header = True

        full_msg += msg.decode('utf-8')

        if len(full_msg) - HEADER_SIZE == msg_len:
            print 'Received: ' + full_msg[HEADER_SIZE:]
            got_msg_header = False
            full_msg = ''
            break
    
    print full_msg 
