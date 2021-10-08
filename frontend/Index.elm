module Index exposing (Model, Msg(..), init, main, update, view)

import Bootstrap.Alert as Alert
import Bootstrap.Button as Button
import Bootstrap.Card as Card
import Bootstrap.Card.Block as Block
import Bootstrap.CDN as CDN
import Bootstrap.Form as Form
import Bootstrap.Form.Input as Input
import Bootstrap.Grid as Grid
import Bootstrap.Grid.Col as Col
import Bootstrap.Grid.Row as Row
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Http
import Json.Decode as Decode exposing (Decoder, int, string)
import Json.Decode.Extra as Decode
import Json.Decode.Pipeline as Pipeline
import Json.Encode as Encode
import Loading exposing (LoaderType(..), defaultConfig, render)



-- MAIN


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }



-- MODEL


type alias Model =
    { username : String
    , password : String
    , state : ModelState
    , serverMessage : String
    }


type ModelState
    = Failure
    | Loading
    | Success
    | Init


type alias Response =
    { message : String
    }


decodeResponse : Decoder Response
decodeResponse =
    Decode.succeed Response
        |> Pipeline.required "message" string


encodeModel : Model -> Encode.Value
encodeModel record =
    Encode.object
        [ ( "username", Encode.string <| record.username )
        , ( "password", Encode.string <| record.password )
        ]


init : () -> ( Model, Cmd Msg )
init _ =
    ( { username = ""
      , password = ""
      , state = Init
      , serverMessage = ""
      }
    , Cmd.none
    )



-- UPDATE


type Msg
    = UpdateUsername String
    | UpdatePassword String
    | SubmitPost
    | GotResponse (Result Http.Error Response)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdateUsername username ->
            ( { model | username = username }, Cmd.none )

        UpdatePassword password ->
            ( { model | password = password }, Cmd.none )

        SubmitPost ->
            ( { model | state = Loading }, sendRequest model )

        GotResponse resp ->
            case resp of
                Ok serverResponse ->
                    ( { model
                        | state = Success
                        , serverMessage = serverResponse.message
                        , username = ""
                        , password = ""
                      }
                    , Cmd.none
                    )

                Err err ->
                    ( { model
                        | state = Failure
                        , serverMessage = errorToString err
                      }
                    , Cmd.none
                    )


errorToString : Http.Error -> String
errorToString error =
    case error of
        Http.BadUrl url ->
            "The URL " ++ url ++ " was invalid"

        Http.Timeout ->
            "Unable to reach the server, try again"

        Http.NetworkError ->
            "Unable to reach the server, check your network connection"

        Http.BadStatus 500 ->
            "The server had a problem, try again later"

        Http.BadStatus 400 ->
            "Verify your information and try again"

        Http.BadStatus 401 ->
            "Authentication failed"

        Http.BadStatus _ ->
            "Unknown error"

        Http.BadBody errorMessage ->
            errorMessage


sendRequest : Model -> Cmd Msg
sendRequest model =
    Http.post
        { body = Http.jsonBody <| encodeModel model
        , url = "/login"
        , expect = Http.expectJson GotResponse decodeResponse
        }


formDisabled : Model -> Bool
formDisabled { state, username, password } =
    if username == "" || password == "" then
        True

    else if List.member state [ Failure, Success, Init ] then
        False

    else
        True



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    Grid.container []
        [ CDN.stylesheet
        , Grid.row [] []
        , Grid.row [ Row.centerXs ]
            [ Grid.col [ Col.xs4 ]
                [ Card.config [ Card.outlineSecondary ]
                    |> Card.headerH4 [] [ text "Go-Web" ]
                    |> Card.block []
                        [ Block.titleH4 [] [ viewResult model ]
                        , Block.custom <|
                            Form.form []
                                [ Form.group []
                                    [ Form.label [ for "username" ] [ text "Username" ]
                                    , Input.text
                                        [ Input.id "username"
                                        , Input.value model.username
                                        , Input.onInput UpdateUsername
                                        , Input.placeholder "Joe"
                                        ]
                                    ]
                                , Form.group []
                                    [ Form.label [ for "password" ] [ text "Password" ]
                                    , Input.password
                                        [ Input.id "password"
                                        , Input.value model.password
                                        , Input.onInput UpdatePassword
                                        , Input.placeholder "secret123"
                                        ]
                                    ]
                                , Button.button
                                    [ Button.primary
                                    , Button.disabled <| formDisabled model
                                    , Button.onClick SubmitPost
                                    ]
                                    [ text "Submit" ]
                                ]
                        ]
                    |> Card.view
                ]
            ]
        ]


viewResult : Model -> Html Msg
viewResult model =
    case model.state of
        Init ->
            div [] []

        Failure ->
            Alert.simpleDanger [] [ text model.serverMessage ]

        Loading ->
            Loading.render
                Circle
                { defaultConfig | color = "#333" }
                Loading.On

        Success ->
            Alert.simpleSuccess [] [ text model.serverMessage ]
