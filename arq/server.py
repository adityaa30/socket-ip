import socket

ADDR = (socket.gethostname(), 1236)

server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(ADDR)
server.listen(5)

WINDOW_SIZE = 4
vals = [False] * WINDOW_SIZE
send_msg = ''

conn, addr = server.accept()
print "Connected to " + str(addr)

while True:
    msg = conn.recv(4096).decode('utf-8')
    orig_msg = msg
    dupl = []
    expect_val = None
    ok = True
    for num in msg.split(','):
        num = int(num)
        if vals[num]:
            dupl.append(num)
        else:
            vals[num] = True

    for i in range(len(vals)):
        if not vals[i]:
            ok = False
            expect_val = i
            break
    
    print(vals)
    if ok:
        vals = [False] * WINDOW_SIZE

        if len(dupl) > 0:
            send_msg = "NCK 0"
        else:
            send_msg = "ACK 0"
    else:
        send_msg = "NCK " + str(expect_val)
            
    conn.send(send_msg.encode('utf-8'))

    