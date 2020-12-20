#!/bin/bash

if [[ $URL == "" ]];then
    echo "You must export the URL environment variable"
    exit 1
fi

TMPFILE=$(mktemp)

if ! which jq 2>&1 > /dev/null;then
    echo "Your system must have 'jq' installed to proceed." >&2
    exit 1
fi

if ! which curl 2>&1 > /dev/null;then
    echo "Your system must have 'curl' installed to proceed." >&2
    exit 2
fi

function ask_for_typecode {
    email_address=$1

    if [[ $email_address == "" ]];then
        echo "ask_for_typecode: you must pass an e-mail address." >&2
        exit 3
    fi

    response=$(
        curl --verbose "$URL/auth/typecodes" \
            -H "channel-name:email" \
            -H "channel-value:$email_address" \
            2> $TMPFILE
    )
    http_status=$(grep "< HTTP/" $TMPFILE)
    if ! echo $http_status | grep -q 200;then
        echo "Something gone wrong when trying to ask for a Typecode." >&2
        echo $http_status >&2
        echo $response >&2
        echo "See $TMPFILE for details" >&2
        exit 4
    fi

    echo "A typecode was sent to your e-mail." >&2
}

function exchange_typecode_for_ticket {
    typecode=$1

    if [[ $typecode == "" ]];then
        echo "exchange_typecode_for_ticket: you must inform a valid Typecode."
        exit 5
    fi

    response=$(
        curl --verbose "$URL/auth/tickets/" \
            -H "typecode:$typecode" \
            2> $TMPFILE
    )
    http_status=$(grep "< HTTP/" $TMPFILE)
    if ! echo $http_status | grep -q 200;then
        echo "Something gone wrong when trying to exchange Typecode for a Ticket" >&2
        echo $http_status >&2
        echo $response >&2
        echo "See $TMPFILE for details" >&2
        echo "TYPECODE=$typecode" >&2
        exit 6
    fi

    ticket=$(echo $response | jq -r .ticket)
    echo $ticket
}

function exchange_ticket_for_access_token {
    ticket=$1

    if [[ $ticket == "" ]];then
        echo "exchange_ticket_for_access_token: you must inform a valid Ticket."
        exit 7
    fi

    response=$(
        curl --verbose "$URL/auth/tokens" \
            -H "ticket:$ticket" \
            2> $TMPFILE
    )
    http_status=$(grep "< HTTP/" $TMPFILE | cut -d'<' -f2-)
    if ! echo $http_status | grep -q 200;then
        echo "Something gone wrong when trying to exchange Ticket for an Access Token" >&2
        echo $http_status >&2
        echo $response >&2
        echo "See $TMPFILE for details" >&2
        echo "TICKET=$ticket" >&2
        exit 6
    fi

    token=$(echo $response | jq -r .access_token)
    echo $token
}

# -----------------------------------
if [[ $TICKET == "" ]];then
    if [[ $TYPECODE == "" ]];then
        # 1- Typecode
        if [[ $EMAIL == "" ]];then
            echo -n "Inform your e-mail address: "
            read EMAIL
        fi

        ask_for_typecode $EMAIL
        echo -n "Type the typecode you received here: "
        read TYPECODE
    fi

    # 2- Ticket
    TICKET=$(exchange_typecode_for_ticket $TYPECODE) || exit $?
fi

# 3- Access Token
token=$(exchange_ticket_for_access_token $TICKET) || exit $?
echo "Your Access Token is:"
echo $token
