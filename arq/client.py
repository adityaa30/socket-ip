import socket, sys, select

ADDR = (socket.gethostname(), 1236)

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(ADDR)

WINDOW_SIZE = 4
curr_val = 0

while True:
    read_streams, _, _ = select.select([sys.stdin, client], [], [])
    for stream in read_streams:
        if stream == client:
            msg = client.recv(4096)
            msg = msg.decode('utf-8')

            print 'Recevived: ', msg
        else:
            msg = str(raw_input())
            orig_msg = msg

            error_msg = 'Packet index should be delimited by "," and in range 0 - ' + str(WINDOW_SIZE - 1)
            
            msg = msg.split(',')
            ok = True
            for num in msg:
                try:
                    num = int(num)
                except:
                    print error_msg
                    ok = False
                    break

                if num >= WINDOW_SIZE or num < 0:
                    print error_msg
                    ok = False
                    break
            
            if ok:
                msg = [int(x) for x in msg]
                client.send(str(orig_msg).encode('utf-8'))
            
