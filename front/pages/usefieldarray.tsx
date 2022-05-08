import { yupResolver } from "@hookform/resolvers/yup";
import { Button, Grid, Theme } from "@mui/material";
import { createStyles, makeStyles } from '@mui/styles';
import { FC } from "react";
import { FormProvider, SubmitHandler, useForm } from "react-hook-form";
import { SchemaOf, string, object } from "yup";

import ReactHookFormTextField from "../src/components/RHookFormText";

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            flexGrow: 1,
            minHeight: '100vh',
        },
    })
);

interface IFormProps {
    PublicKey: string
    PrivateKey: string
    BlockchainAddress: string
}

const formSchema: SchemaOf<IFormProps> = object({
    PublicKey: string().required('入力必須です'),
    PrivateKey: string().required('入力必須です'),
    BlockchainAddress: string().required('入力必須です'),
})

const FieldArrayForm: FC = () => {
    const classes = useStyles()

    const methods = useForm<IFormProps>({
        resolver: yupResolver(formSchema)
    })

    const submit: SubmitHandler<IFormProps> = async (data: IFormProps) => {
        console.log('data submitted', data)
    }

    return (
        <div>
            <Grid container>
                <FormProvider {...methods}>
                    <form onSubmit={methods.handleSubmit(submit)}>
                        <Grid item>
                            <ReactHookFormTextField label="Public Key" name="PublicKey" />
                        </Grid>
                        <Grid item>
                            <ReactHookFormTextField label="Private Key" name="PrivateKey" />
                        </Grid>
                        <Grid item>
                            <ReactHookFormTextField label="Blockchain Address" name="BlockchainAddress" />
                        </Grid>
                        <Grid item>
                            <Button
                                type="submit"
                                variant="contained"
                                color="primary"
                            >
                                送信
                            </Button>
                        </Grid>
                    </form>
                </FormProvider>
            </Grid>
        </div>
    )
}

export default FieldArrayForm
